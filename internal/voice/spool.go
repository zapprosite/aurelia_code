package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const jobManifestName = "job.json"

type Job struct {
	ID            string    `json:"id"`
	Source        string    `json:"source"`
	UserID        int64     `json:"user_id"`
	ChatID        int64     `json:"chat_id"`
	AudioFile     string    `json:"audio_file"`
	RequiresAudio bool      `json:"requires_audio"`
	CreatedAt     time.Time `json:"created_at"`
}

type Result struct {
	Status      string    `json:"status"`
	Transcript  string    `json:"transcript,omitempty"`
	Error       string    `json:"error,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}

type ClaimedJob struct {
	Dir      string
	AudioDir string
	Job      Job
}

type Spool struct {
	root       string
	inboxDir   string
	workingDir string
	doneDir    string
	failedDir  string
}

func NewSpool(root string) (*Spool, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, fmt.Errorf("voice spool root is required")
	}
	spool := &Spool{
		root:       root,
		inboxDir:   filepath.Join(root, "inbox"),
		workingDir: filepath.Join(root, "processing"),
		doneDir:    filepath.Join(root, "done"),
		failedDir:  filepath.Join(root, "failed"),
	}
	for _, dir := range []string{spool.root, spool.inboxDir, spool.workingDir, spool.doneDir, spool.failedDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("create voice spool dir %q: %w", dir, err)
		}
	}
	return spool, nil
}

func (s *Spool) Root() string { return s.root }

func (s *Spool) EnqueueAudioFile(job Job, sourcePath string) (Job, error) {
	if s == nil {
		return Job{}, fmt.Errorf("voice spool is nil")
	}
	if strings.TrimSpace(sourcePath) == "" {
		return Job{}, fmt.Errorf("source audio path is required")
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return Job{}, fmt.Errorf("stat source audio: %w", err)
	}
	if sourceInfo.IsDir() {
		return Job{}, fmt.Errorf("source audio path must be a file")
	}

	if job.ID == "" {
		job.ID = fmt.Sprintf("%d-%s", time.Now().UTC().UnixNano(), uuid.NewString())
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	}
	if job.Source == "" {
		job.Source = "unknown"
	}

	tempDir, err := os.MkdirTemp(s.inboxDir, "job-*")
	if err != nil {
		return Job{}, fmt.Errorf("create temp voice job dir: %w", err)
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(tempDir)
		}
	}()

	ext := filepath.Ext(sourcePath)
	if ext == "" {
		ext = ".bin"
	}
	job.AudioFile = "audio" + ext
	destAudioPath := filepath.Join(tempDir, job.AudioFile)
	if err := copyFile(sourcePath, destAudioPath); err != nil {
		return Job{}, err
	}
	if err := writeJSONFile(filepath.Join(tempDir, jobManifestName), job); err != nil {
		return Job{}, err
	}

	finalDir := filepath.Join(s.inboxDir, job.ID)
	if err := os.Rename(tempDir, finalDir); err != nil {
		return Job{}, fmt.Errorf("move voice job into inbox: %w", err)
	}
	success = true
	return job, nil
}

func (s *Spool) ClaimOldest(ctx context.Context) (*ClaimedJob, error) {
	if s == nil {
		return nil, fmt.Errorf("voice spool is nil")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	entries, err := os.ReadDir(s.inboxDir)
	if err != nil {
		return nil, fmt.Errorf("read voice inbox: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		inboxPath := filepath.Join(s.inboxDir, name)
		workingPath := filepath.Join(s.workingDir, name)
		if err := os.Rename(inboxPath, workingPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("claim voice job %q: %w", name, err)
		}

		job, err := readJobManifest(filepath.Join(workingPath, jobManifestName))
		if err != nil {
			return nil, err
		}
		return &ClaimedJob{
			Dir:      workingPath,
			AudioDir: filepath.Join(workingPath, job.AudioFile),
			Job:      job,
		}, nil
	}
	return nil, nil
}

func (s *Spool) Complete(claimed *ClaimedJob, transcript string) error {
	return s.finish(claimed, s.doneDir, Result{
		Status:      "done",
		Transcript:  strings.TrimSpace(transcript),
		ProcessedAt: time.Now().UTC(),
	})
}

func (s *Spool) Fail(claimed *ClaimedJob, err error) error {
	message := ""
	if err != nil {
		message = strings.TrimSpace(err.Error())
	}
	return s.finish(claimed, s.failedDir, Result{
		Status:      "failed",
		Error:       message,
		ProcessedAt: time.Now().UTC(),
	})
}

func (s *Spool) QueueDepth() (int, error) {
	if s == nil {
		return 0, fmt.Errorf("voice spool is nil")
	}
	depth := 0
	for _, dir := range []string{s.inboxDir, s.workingDir} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return 0, fmt.Errorf("read voice queue depth from %q: %w", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				depth++
			}
		}
	}
	return depth, nil
}

func (s *Spool) finish(claimed *ClaimedJob, targetDir string, result Result) error {
	if s == nil {
		return fmt.Errorf("voice spool is nil")
	}
	if claimed == nil {
		return fmt.Errorf("claimed voice job is nil")
	}
	if err := writeJSONFile(filepath.Join(claimed.Dir, "result.json"), result); err != nil {
		return err
	}
	destination := filepath.Join(targetDir, claimed.Job.ID)
	if err := os.RemoveAll(destination); err != nil {
		return fmt.Errorf("cleanup previous voice job result: %w", err)
	}
	if err := os.Rename(claimed.Dir, destination); err != nil {
		return fmt.Errorf("move voice job result: %w", err)
	}
	return nil
}

func readJobManifest(path string) (Job, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Job{}, fmt.Errorf("read voice job manifest: %w", err)
	}
	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return Job{}, fmt.Errorf("decode voice job manifest: %w", err)
	}
	if job.ID == "" || job.AudioFile == "" {
		return Job{}, fmt.Errorf("voice job manifest missing required fields")
	}
	return job, nil
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode json file %q: %w", path, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write json file %q: %w", path, err)
	}
	return nil
}

func copyFile(sourcePath, destPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source audio file: %w", err)
	}
	defer func() { _ = source.Close() }()

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create destination audio file: %w", err)
	}
	defer func() { _ = dest.Close() }()

	if _, err := io.Copy(dest, source); err != nil {
		return fmt.Errorf("copy audio file into spool: %w", err)
	}
	if err := dest.Close(); err != nil {
		return fmt.Errorf("close destination audio file: %w", err)
	}
	return nil
}

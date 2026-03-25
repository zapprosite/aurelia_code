import * as React from "react";
import { X } from "lucide-react";

interface PersonaTemplate {
  id: string;
  name: string;
  description: string;
  icon: string;
  color: string;
}

interface CreateBotModalProps {
  onClose: () => void;
  onCreated: () => void;
}

export function CreateBotModal({ onClose, onCreated }: CreateBotModalProps) {
  const [personas, setPersonas] = React.useState<PersonaTemplate[]>([]);
  const [form, setForm] = React.useState({
    id: "",
    name: "",
    token: "",
    persona_id: "",
    focus_area: "",
    allowed_user_ids_raw: "",
  });
  const [error, setError] = React.useState<string | null>(null);
  const [loading, setLoading] = React.useState(false);

  React.useEffect(() => {
    fetch("/api/personas")
      .then((r) => r.json())
      .then(setPersonas)
      .catch(() => {});
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const allowed_user_ids = form.allowed_user_ids_raw
      .split(",")
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));

    try {
      const res = await fetch("/api/bots/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          id: form.id,
          name: form.name,
          token: form.token,
          persona_id: form.persona_id,
          focus_area: form.focus_area,
          allowed_user_ids: allowed_user_ids.length > 0 ? allowed_user_ids : null,
          enabled: true,
        }),
      });
      if (!res.ok) {
        const msg = await res.text();
        setError(msg || `Erro ${res.status}`);
        return;
      }
      onCreated();
      onClose();
    } catch (err) {
      setError("Falha de rede ao criar bot");
    } finally {
      setLoading(false);
    }
  };

  const field = (
    label: string,
    key: keyof typeof form,
    placeholder: string,
    required = false,
    hint?: string
  ) => (
    <div className="flex flex-col gap-1.5">
      <label className="text-[11px] font-mono uppercase tracking-widest text-white/50">
        {label} {required && <span className="text-red-400">*</span>}
      </label>
      <input
        className="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white/90 placeholder-white/20 focus:outline-none focus:border-white/30 transition-colors"
        placeholder={placeholder}
        value={form[key]}
        required={required}
        onChange={(e) => setForm((f) => ({ ...f, [key]: e.target.value }))}
      />
      {hint && <span className="text-[10px] text-white/30">{hint}</span>}
    </div>
  );

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="relative w-full max-w-lg rounded-2xl border border-white/10 bg-[#161616]/95 backdrop-blur-xl p-8 shadow-2xl">
        {/* Close */}
        <button
          className="absolute top-5 right-5 text-white/30 hover:text-white/70 transition-colors"
          onClick={onClose}
        >
          <X className="w-5 h-5" />
        </button>

        <h2 className="text-base font-semibold text-white/90 mb-1">Adicionar Bot</h2>
        <p className="text-xs text-white/40 mb-6">Configure um novo membro do time Telegram.</p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-2 gap-4">
            {field("ID", "id", "hvac-sales", true, "Identificador único (sem espaços)")}
            {field("Nome", "name", "Bot de Vendas", true)}
          </div>

          {field("Token Telegram", "token", "123456789:AAF...", true, "Token do @BotFather")}

          {/* Persona dropdown */}
          <div className="flex flex-col gap-1.5">
            <label className="text-[11px] font-mono uppercase tracking-widest text-white/50">
              Persona
            </label>
            <select
              className="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white/90 focus:outline-none focus:border-white/30 transition-colors"
              value={form.persona_id}
              onChange={(e) => setForm((f) => ({ ...f, persona_id: e.target.value }))}
            >
              <option value="">— Selecione uma persona —</option>
              {personas.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>

          {field("Área de Foco", "focus_area", "Funil DAIKIN VRV SP, specs, preços")}
          {field(
            "Usuários Permitidos",
            "allowed_user_ids_raw",
            "123456789, 987654321",
            false,
            "IDs Telegram separados por vírgula. Vazio = usa os globais."
          )}

          {error && (
            <div className="rounded-lg bg-red-500/10 border border-red-500/20 px-4 py-2 text-sm text-red-400">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="mt-2 w-full rounded-lg bg-primary/80 hover:bg-primary text-white font-semibold py-2.5 text-sm transition-colors disabled:opacity-50"
          >
            {loading ? "Criando..." : "Criar Bot"}
          </button>
        </form>
      </div>
    </div>
  );
}

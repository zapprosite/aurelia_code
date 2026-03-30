package telegram

const (
	unsupportedDocumentMessage = "⚠️ **Formato não suportado**\n\n" +
		"Para garantir a melhor análise, no momento consigo processar os seguintes formatos:\n" +
		"- Documentos Markdown (`.md`)\n" +
		"- Arquivos PDF (`.pdf`)\n" +
		"- Mensagens de áudio e voz\n" +
		"- Imagens e fotos (incluindo álbuns)\n"

	downloadFailureMessage = "❌ **Falha no processamento**\n\n" +
		"Não foi possível baixar o arquivo enviado. Por favor, tente encaminhar novamente."

	audioNotConfiguredMessage = "⚠️ **Módulo de Áudio não configurado**\n\n" +
		"O serviço de transcrição ainda não foi ativado.\n\n" +
		"Para habilitar, configure a `groq_api_key` em suas definições de sistema."

	audioProcessingFailureMessage = "🎙️ **Serviço de Voz**\n\n" +
		"Desculpe, não consegui converter seu áudio em texto com a clareza necessária agora. Poderia repetir ou enviar uma mensagem de texto?"

	emptyAudioMessage = "⚠️ **Áudio sem conteúdo**\n\n" +
		"O arquivo de áudio parece estar vazio ou inaudível. Poderia verificar e reenviar?"

	alreadyConfiguredMessage = "**Aurélia online.** Manda sua tarefa."

	bootstrapWelcomeMessage = "**Aurélia** — selecione o perfil:"

	bootstrapFailureMessage = "Erro ao criar identidade. Tenta `/start` novamente."

	bootstrapProfileMessage = "Perfil configurado. Qual é o seu nome e como prefere ser tratado?"

	bootstrapSuccessMessage = "Pronta. Manda sua tarefa — texto, áudio, imagem ou comando."
)

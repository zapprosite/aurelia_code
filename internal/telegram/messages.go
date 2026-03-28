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

	alreadyConfiguredMessage = "✨ **Sistema Operacional**\n\n" +
		"A Aurélia já está configurada e pronta para auxiliá-lo. Como posso ser útil hoje?"

	bootstrapWelcomeMessage = "💎 **Boas-vindas à Aurélia**\n\n" +
		"Sou sua nova assistente de inteligência artificial de elite.\n\n" +
		"Para começarmos, selecione o perfil de atuação que melhor atende às suas necessidades atuais:"

	bootstrapFailureMessage = "❌ **Falha na Inicialização**\n\n" +
		"Ocorreu um erro técnico ao criar seus arquivos de identidade. Por favor, contate o administrador."

	bootstrapProfileMessage = "✅ **Perfil Inicial Configurado**\n\n" +
		"Excelente escolha. Agora, por favor, informe seu nome e como você gostaria que eu me dirigisse a você (formal ou informal).\n\n" +
		"Exemplo: `Sou o Rafael e prefiro uma comunicação direta e formal.`"

	bootstrapSuccessMessage = "🎯 **Configuração Concluída**\n\n" +
		"Seus protocolos de identidade foram salvos com sucesso.\n\n" +
		"Estou pronta para atuar. Sinta-se à vontade para enviar textos, comandos, áudios ou imagens para análise."
)

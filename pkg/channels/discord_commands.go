package channels

import (
	"github.com/bwmarrin/discordgo"
)

var DiscordCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Iniciar el bot y ver bienvenida",
	},
	{
		Name:        "help",
		Description: "Mostrar mensaje de ayuda con commandos disponibles",
	},
	{
		Name:        "model",
		Description: "Cambiar o listar modelos de IA",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "model",
				Description: "Nombre del modelo o provider (ej: openai/gpt-5.4)",
				Required:    false,
			},
		},
	},
	{
		Name:        "models",
		Description: "Listar todos los modelos disponibles",
	},
	{
		Name:        "bundle_approve",
		Description: "Aprobar un lote de post+imagen",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ID del lote (ej: 20260302_161740_yiia22)",
				Required:    true,
			},
		},
	},
	{
		Name:        "bundle_regen",
		Description: "Regenerar un lote completo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ID del lote",
				Required:    true,
			},
		},
	},
	{
		Name:        "bundle_edit",
		Description: "Editar el texto de un lote",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ID del lote",
				Required:    true,
			},
		},
	},
	{
		Name:        "bundle_publish",
		Description: "Publicar un lote aprobado en redes sociales",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ID del lote",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "platforms",
				Description: "Plataformas (ej: facebook,twitter,discord)",
				Required:    false,
			},
		},
	},
	{
		Name:        "bundle_cancel",
		Description: "Cancelar un lote y descartarlo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ID del lote",
				Required:    true,
			},
		},
	},
	{
		Name:        "show",
		Description: "Ver configuración actual",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "target",
				Description: "Objetivo: 'model' o 'channel'",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Model", Value: "model"},
					{Name: "Channel", Value: "channel"},
				},
			},
		},
	},
	{
		Name:        "list",
		Description: "Listar opciones disponibles",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "target",
				Description: "Objetivo: 'models' o 'channels'",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Models", Value: "models"},
					{Name: "Channels", Value: "channels"},
				},
			},
		},
	},
	{
		Name:        "status",
		Description: "Ver estado del contexto y tokens",
	},
	{
		Name:        "disable_sentinel",
		Description: "Desactivar sentinel temporalmente (como antivirus)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "duration",
				Description: "Duración de la desactivación",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "5 minutos", Value: "5m"},
					{Name: "15 minutos", Value: "15m"},
					{Name: "1 hora", Value: "1h"},
				},
			},
		},
	},
	{
		Name:        "activate_sentinel",
		Description: "Activar sentinel inmediatamente",
	},
	{
		Name:        "sentinel_status",
		Description: "Ver estado actual del sentinel",
	},
	{
		Name:        "restrict_to_workspace",
		Description: "Controlar restricción de acceso a archivos del sistema",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Acción a ejecutar",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "activate — Solo workspace (SEGURO)", Value: "activate"},
					{Name: "deactivate — Acceso total (PELIGROSO)", Value: "deactivate"},
					{Name: "status — Ver estado actual", Value: "status"},
				},
			},
		},
	},
}

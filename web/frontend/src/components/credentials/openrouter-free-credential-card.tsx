import React from "react"
import { IconKey, IconLoader2 } from "@tabler/icons-react"
import { useTranslation } from "react-i18next"

import type { OAuthProviderStatus } from "@/api/oauth"
import { Button } from "@/components/ui/button"

import { CredentialCard } from "./credential-card"

interface OpenRouterFreeCredentialCardProps {
  status?: OAuthProviderStatus
  activeAction: string
  onAskLogout: () => void
}

export function OpenRouterFreeCredentialCard({
  status,
  activeAction,
  onAskLogout,
}: OpenRouterFreeCredentialCardProps) {
  const { t } = useTranslation()
  const actionBusy = activeAction !== ""
  const [token, setToken] = React.useState("")
  const [isSaving, setIsSaving] = React.useState(false)

  const handleSaveToken = async () => {
    if (!token || actionBusy) return

    setIsSaving(true)
    try {
      const response = await fetch('/api/oauth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: 'openrouter-free',
          method: 'token',
          token: token.trim()
        })
      })

      if (!response.ok) {
        const error = await response.text()
        throw new Error(error || 'Failed to save API key')
      }

      // Success - reload page to show updated status
      window.location.reload()
    } catch (error) {
      console.error('Error saving OpenRouter Free API key:', error)
      alert('Failed to save API key: ' + (error as Error).message)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <CredentialCard
      title={
        <span className="inline-flex items-center gap-2">
          <span className="border-muted inline-flex size-6 items-center justify-center rounded-full border">
            <span className="text-xs font-bold leading-none">OR</span>
          </span>
          <span>OpenRouter Free</span>
        </span>
      }
      description="Free tier: Auto-routes to best available free model"
      status={status?.status ?? "not_logged_in"}
      authMethod={status?.auth_method}
      details={
        <div className="space-y-1">
          {!status?.logged_in && (
            <p className="text-muted-foreground text-xs">
              Get your API key at:{" "}
              <a
                href="https://openrouter.ai/keys"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline"
              >
                openrouter.ai/keys
              </a>
              {" • "}
              <a
                href="https://openrouter.ai/collections/free-models"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline"
              >
                Free models
              </a>
            </p>
          )}
        </div>
      }
      actions={
        <div className="border-muted flex h-[120px] flex-col justify-center rounded-lg border p-3">
          {!status?.logged_in ? (
            <div className="space-y-2">
              <input
                type="password"
                placeholder="sk-or-v1-..."
                value={token}
                onChange={(e) => setToken(e.target.value)}
                className="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm focus-visible:outline-none focus-visible:ring-1"
              />
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  disabled={actionBusy || isSaving || !token}
                  onClick={handleSaveToken}
                >
                  {isSaving && (
                    <IconLoader2 className="size-4 animate-spin" />
                  )}
                  <IconKey className="size-4" />
                  {isSaving ? 'Saving...' : 'Save API Key'}
                </Button>
              </div>
              <p className="text-muted-foreground text-xs">
                ✅ 100% Free tier • No model selection needed • Auto-routes
              </p>
            </div>
          ) : (
            <div className="flex flex-wrap items-center gap-2">
              <Button
                size="sm"
                variant="outline"
                disabled={actionBusy}
                onClick={onAskLogout}
              >
                <IconKey className="size-4" />
                API Key configured
              </Button>
            </div>
          )}
        </div>
      }
      footer={
        status?.logged_in ? (
          <Button
            variant="ghost"
            size="sm"
            disabled={actionBusy}
            onClick={onAskLogout}
            className="text-destructive hover:bg-destructive/10 hover:text-destructive"
          >
            {activeAction === "openrouter-free:logout" && (
              <IconLoader2 className="size-4 animate-spin" />
            )}
            {t("credentials.actions.logout")}
          </Button>
        ) : null
      }
    />
  )
}

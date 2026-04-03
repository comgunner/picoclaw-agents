import React from "react"
import { IconKey, IconLoader2 } from "@tabler/icons-react"
import { useTranslation } from "react-i18next"

import type { OAuthProviderStatus } from "@/api/oauth"
import { Button } from "@/components/ui/button"

import { CredentialCard } from "./credential-card"

interface QwenCredentialCardProps {
  status?: OAuthProviderStatus
  activeAction: string
  onAskLogout: () => void
}

export function QwenCredentialCard({
  status,
  activeAction,
  onAskLogout,
}: QwenCredentialCardProps) {
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
          provider: 'qwen',
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
      console.error('Error saving Qwen API key:', error)
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
            <span className="text-xs font-bold leading-none">Q</span>
          </span>
          <span>Qwen Portal</span>
        </span>
      }
      description="API key authentication for Qwen Portal (DashScope)"
      status={status?.status ?? "not_logged_in"}
      authMethod={status?.auth_method}
      details={
        <div className="space-y-1">
          {status?.account_id && (
            <p>
              {t("credentials.labels.account")}: {status.account_id}
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
                placeholder="sk-..."
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
            {activeAction === "qwen:logout" && (
              <IconLoader2 className="size-4 animate-spin" />
            )}
            {t("credentials.actions.logout")}
          </Button>
        ) : null
      }
    />
  )
}

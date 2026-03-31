type HeaderProps = {
  showConnectionStatus: boolean
  isConnected: boolean
  backendStatus: 'checking' | 'healthy' | 'unhealthy'
}

export function Header({ showConnectionStatus, isConnected, backendStatus }: HeaderProps) {
  const backendBadgeClass =
    backendStatus === 'healthy'
      ? 'border-level-INFO/40 bg-level-INFO/10 text-level-INFO'
      : backendStatus === 'unhealthy'
        ? 'border-level-ERROR/40 bg-level-ERROR/10 text-level-ERROR'
        : 'border-border bg-bg/50 text-text-muted'

  return (
    <header className="flex items-center justify-between border-b border-border bg-surface/80 px-4 py-4 backdrop-blur sm:px-6">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-border bg-bg/80 text-xl text-text-primary shadow-glow">
          ☁
        </div>
        <div>
          <p className="font-sans text-lg font-semibold tracking-tight text-text-primary">LogFlow</p>
          <p className="font-sans text-xs uppercase tracking-[0.24em] text-text-muted">System log monitor</p>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <div className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 font-sans text-xs font-semibold tracking-[0.16em] ${backendBadgeClass}`}>
          <span className={`text-[10px] ${backendStatus === 'healthy' ? 'animate-pulse-dot' : ''}`}>●</span>
          {backendStatus === 'healthy'
            ? 'GATEWAY OK'
            : backendStatus === 'unhealthy'
              ? 'GATEWAY DOWN'
              : 'CHECKING'}
        </div>

        {showConnectionStatus ? (
          <div
            className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 font-sans text-xs font-semibold tracking-[0.2em] ${
              isConnected
                ? 'border-level-INFO/40 bg-level-INFO/10 text-level-INFO'
                : 'border-level-ERROR/40 bg-level-ERROR/10 text-level-ERROR'
            }`}
          >
          <span className={`text-[10px] ${isConnected ? 'animate-pulse-dot' : ''}`}>●</span>
          {isConnected ? 'LIVE' : 'DISCONNECTED'}
          </div>
        ) : null}
      </div>
    </header>
  )
}

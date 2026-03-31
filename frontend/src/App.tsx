import { useState } from 'react'

import { Header } from './components/Header'
import { LiveLogsTab } from './components/LiveLogs/LiveLogsTab'
import { SearchTab } from './components/Search/SearchTab'
import { TabBar } from './components/TabBar'
import { useBackendHealth } from './hooks/useBackendHealth'

type TabKey = 'live' | 'search'

function App() {
  const [activeTab, setActiveTab] = useState<TabKey>('live')
  const [isConnected, setIsConnected] = useState(false)
  const { status, isBackendReachable, checkedAtLeastOnce, backendLabel } = useBackendHealth()

  return (
    <div className="min-h-screen bg-bg text-text-primary">
      <div className="mx-auto min-h-screen max-w-[1600px] border-x border-border bg-[radial-gradient(circle_at_top,_rgba(88,166,255,0.14),_transparent_34%),linear-gradient(180deg,_rgba(22,27,34,0.98),_rgba(13,17,23,1))]">
        <Header
          showConnectionStatus={activeTab === 'live'}
          isConnected={checkedAtLeastOnce && isBackendReachable && isConnected}
          backendStatus={status}
        />
        <TabBar activeTab={activeTab} onChange={setActiveTab} />
        <main className="px-4 py-6 sm:px-6">
          {checkedAtLeastOnce && !isBackendReachable ? (
            <div className="mb-4 rounded-2xl border border-level-ERROR/40 bg-level-ERROR/10 px-4 py-3 font-sans text-sm text-level-ERROR">
              Backend is down or unreachable at <span className="font-mono">{backendLabel}</span>.
            </div>
          ) : null}
          {activeTab === 'live' ? (
            <LiveLogsTab onConnectionChange={setIsConnected} />
          ) : (
            <SearchTab />
          )}
        </main>
      </div>
    </div>
  )
}

export default App

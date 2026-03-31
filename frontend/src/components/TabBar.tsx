type TabKey = 'live' | 'search'

type TabBarProps = {
  activeTab: TabKey
  onChange: (tab: TabKey) => void
}

export function TabBar({ activeTab, onChange }: TabBarProps) {
  const tabs: Array<{ key: TabKey; label: string }> = [
    { key: 'live', label: 'Live Logs' },
    { key: 'search', label: 'Search' },
  ]

  return (
    <div className="border-b border-border px-4 sm:px-6">
      <div className="flex gap-2 py-4">
        {tabs.map((tab) => {
          const active = activeTab === tab.key
          return (
            <button
              key={tab.key}
              type="button"
              onClick={() => onChange(tab.key)}
              className={`rounded-full border px-4 py-2 font-sans text-sm font-medium transition ${
                active
                  ? 'border-[#58a6ff]/50 bg-[#58a6ff]/12 text-text-primary shadow-glow'
                  : 'border-border bg-surface text-text-muted hover:border-[#58a6ff]/30 hover:text-text-primary'
              }`}
            >
              {tab.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}

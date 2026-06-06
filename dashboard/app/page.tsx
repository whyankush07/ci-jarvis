import { getRuns } from "../lib/api";
import { RunCard } from "../components/dashboard/RunCard";

export default async function DashboardPage() {
  const data = await getRuns();
  const runs = data.runs;
  console.log(data.runs)

  return (
    <div className="max-w-5xl mx-auto space-y-8">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
          Jarvis
        </h1>
        <p className="text-zinc-500 dark:text-zinc-400">
          Agentic CI/CD Pipeline Monitoring
        </p>
      </header>

      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-zinc-800 dark:text-zinc-200">
            Recent Runs
          </h2>
          <div className="h-px flex-1 bg-zinc-100 dark:bg-zinc-800 mx-4" />
          <span className="text-xs font-medium text-zinc-400 uppercase tracking-wider">
            Live Updates
          </span>
        </div>

        {runs && runs.length === 0 ? (
          <div className="py-12 text-center border-2 border-dashed border-zinc-100 dark:border-zinc-800 rounded-2xl">
            <p className="text-zinc-400 dark:text-zinc-500">No active runs found</p>
          </div>
        ) : (
          <div className="grid gap-3">
            {runs && runs.map((run) => (
              <RunCard key={run.id} run={run} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}

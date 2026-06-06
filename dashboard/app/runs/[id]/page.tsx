import { getRun } from "@/lib/api";
import { Badge } from "@/components/ui/Badge";
import { Card, CardContent, CardHeader } from "@/components/ui/Card";
import { formatDate, getStatusVariant } from "@/lib/utils";
import Link from "next/link";
import { 
  PlannerOutput, 
  CoderOutput, 
  ReviewerOutput 
} from "@/types";

export default async function RunDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const resolvedParams = await params;
  const run = await getRun(resolvedParams.id);
  const plannerData = run.plan as PlannerOutput | null;
  const coderData = run.metadata?.coder_output as CoderOutput | null;
  const reviewerData = run.metadata?.reviewer_output as ReviewerOutput | null;
  const ragContext = run.metadata?.RelatedContext as string | null;

  return (
    <div className="max-w-5xl mx-auto space-y-10">
      <header className="space-y-4">
        <Link 
          href="/" 
          className="text-sm text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200 transition-colors flex items-center gap-1"
        >
          ← Back to Dashboard
        </Link>
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div className="space-y-1">
            <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
              {run.repo_url.split("//").pop()?.replace(".git", "")}
            </h1>
            <p className="text-zinc-500 dark:text-zinc-400 font-mono text-sm">
              Run ID: {run.id}
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Badge variant={getStatusVariant(run.status)}>
              {run.status.replace("_", " ")}
            </Badge>
            <span className="text-sm text-zinc-400">
              {formatDate(run.created_at)}
            </span>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-1 gap-8">
        <section className="space-y-4">
          <div className="flex items-center gap-4">
            <div className="w-8 h-8 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-xs font-bold">1</div>
            <h2 className="text-lg font-semibold">Strategic Planning</h2>
          </div>
          {plannerData ? (
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <h3 className="font-medium">Execution Strategy</h3>
                  <Badge variant={plannerData.priority === "high" ? "error" : plannerData.priority === "medium" ? "warning" : "default"}>
                    {plannerData.priority} priority
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-6">
                <div>
                  <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Summary</h4>
                  <p className="text-zinc-700 dark:text-zinc-300">{plannerData.summary}</p>
                </div>
                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Planned Steps</h4>
                    <ul className="space-y-2">
                      {plannerData.steps.map((step, i) => (
                        <li key={i} className="text-sm flex gap-2">
                          <span className="text-zinc-400">{i + 1}.</span>
                          <span className="text-zinc-600 dark:text-zinc-400">{step}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                  <div>
                    <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Identified Risks</h4>
                    <ul className="space-y-2">
                      {plannerData.risks.map((risk, i) => (
                        <li key={i} className="text-sm flex gap-2">
                          <span className="text-rose-500">⚠</span>
                          <span className="text-zinc-600 dark:text-zinc-400">{risk}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </CardContent>
            </Card>
          ) : (
            <div className="p-8 text-center bg-zinc-50 dark:bg-zinc-900/50 rounded-xl border border-zinc-100 dark:border-zinc-800 italic text-zinc-400">
              Planning phase in progress...
            </div>
          )}
        </section>

        <section className="space-y-4">
          <div className="flex items-center gap-4">
            <div className="w-8 h-8 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-xs font-bold">2</div>
            <h2 className="text-lg font-semibold">Code Generation</h2>
          </div>
          {coderData ? (
            <div className="space-y-4">
              <Card>
                <CardContent className="py-6">
                  <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Implementation Explanation</h4>
                  <p className="text-zinc-700 dark:text-zinc-300">{coderData.explanation}</p>
                </CardContent>
              </Card>
              {coderData.changes.map((change, i) => (
                <Card key={i} className="border-l-4 border-l-emerald-500">
                  <CardHeader className="py-3 bg-zinc-50 dark:bg-zinc-800/50">
                    <div className="flex items-center justify-between">
                      <code className="text-xs font-mono">{change.path}</code>
                      <Badge variant="success">{change.action}</Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="p-0">
                    <pre className="p-4 text-xs font-mono overflow-x-auto bg-black text-zinc-300">
                      <code>{change.content}</code>
                    </pre>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <div className="p-8 text-center bg-zinc-50 dark:bg-zinc-900/50 rounded-xl border border-zinc-100 dark:border-zinc-800 italic text-zinc-400">
              {run.status === "planning" ? "Waiting for planning to finish..." : "Coding phase in progress..."}
            </div>
          )}
        </section>

        <section className="space-y-4">
          <div className="flex items-center gap-4">
            <div className="w-8 h-8 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-xs font-bold">3</div>
            <h2 className="text-lg font-semibold">AI Peer Review</h2>
          </div>
          {reviewerData ? (
            <Card className={reviewerData.approved ? "border-l-4 border-l-emerald-500" : "border-l-4 border-l-rose-500"}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <h3 className="font-medium">Review Results</h3>
                  <Badge variant={reviewerData.approved ? "success" : "error"}>
                    {reviewerData.approved ? "Approved" : "Changes Requested"}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-6">
                <div>
                  <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Review Summary</h4>
                  <p className="text-zinc-700 dark:text-zinc-300">{reviewerData.summary}</p>
                </div>
                {reviewerData.comments.length > 0 && (
                  <div>
                    <h4 className="text-xs font-bold uppercase tracking-widest text-zinc-400 mb-2">Reviewer Comments</h4>
                    <div className="space-y-3">
                      {reviewerData.comments.map((comment, i) => (
                        <div key={i} className="p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-100 dark:border-zinc-800">
                          <div className="flex items-center justify-between mb-1">
                            <code className="text-[10px] font-mono text-zinc-500">{comment.file}:{comment.line}</code>
                            <Badge variant={comment.level === "error" ? "error" : comment.level === "warning" ? "warning" : "info"}>
                              {comment.level}
                            </Badge>
                          </div>
                          <p className="text-sm text-zinc-600 dark:text-zinc-400">{comment.comment}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          ) : (
            <div className="p-8 text-center bg-zinc-50 dark:bg-zinc-900/50 rounded-xl border border-zinc-100 dark:border-zinc-800 italic text-zinc-400">
              Waiting for review phase...
            </div>
          )}
        </section>

        {ragContext && (
          <section className="space-y-4">
            <div className="flex items-center gap-4">
              <div className="w-8 h-8 rounded-full bg-sky-50 dark:bg-sky-900/20 flex items-center justify-center text-sky-500">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m21 21-6-6m6 6v-4.8m0 4.8h-4.8M3 3l6 6m-6-6v4.8m0-4.8h4.8m0 13.2V21m0 0H3m0 0 6-6m12-12V3m0 0h-4.8M21 3l-6 6"/></svg>
              </div>
              <h2 className="text-lg font-semibold">Retrieved Repository Context</h2>
            </div>
            <Card className="bg-sky-50/30 dark:bg-sky-900/5 border-sky-100 dark:border-sky-900/20">
              <CardContent className="py-6">
                <div className="prose prose-sm dark:prose-invert max-w-none">
                  <pre className="whitespace-pre-wrap font-mono text-[10px] leading-relaxed text-sky-900/70 dark:text-sky-400/70 bg-transparent p-0">
                    {ragContext}
                  </pre>
                </div>
              </CardContent>
            </Card>
          </section>
        )}
      </div>
    </div>
  );
}

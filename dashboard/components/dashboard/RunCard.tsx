import Link from "next/link";
import { Run } from "../../types";
import { Badge } from "../ui/Badge";
import { Card, CardContent } from "../ui/Card";
import { formatDate, getStatusVariant } from "../../lib/utils";

interface RunCardProps {
  run: Run;
}

export const RunCard = ({ run }: RunCardProps) => {
  return (
    <Link href={`/runs/${run.id}`}>
      <Card className="hover:border-zinc-300 dark:hover:border-zinc-700 transition-colors group">
        <CardContent className="py-4">
          <div className="flex items-center justify-between gap-4">
            <div className="flex flex-col min-w-0 flex-1">
              <div className="flex items-center gap-2 mb-1">
                <span className="font-semibold truncate text-zinc-900 dark:text-zinc-100">
                  {run.repo_url.split("//").pop()?.replace(".git", "") || "Repository"}
                </span>
                <Badge variant={getStatusVariant(run.status)}>
                  {run.status.replace("_", " ")}
                </Badge>
              </div>
              <p className="text-sm text-zinc-500 dark:text-zinc-400 truncate">
                {run.current_step || "Initializing..."}
              </p>
            </div>
            <div className="text-right flex-shrink-0">
              <p className="text-xs text-zinc-400 dark:text-zinc-500 mb-1">
                {formatDate(run.created_at)}
              </p>
              <div className="text-xs font-mono text-zinc-300 dark:text-zinc-700 group-hover:text-zinc-400 dark:group-hover:text-zinc-600 transition-colors">
                {run.id.slice(0, 8)}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
};

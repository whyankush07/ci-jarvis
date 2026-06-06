export type RunStatus = 
  | "pending" 
  | "planning" 
  | "coding" 
  | "reviewing" 
  | "executing" 
  | "awaiting_approval" 
  | "completed" 
  | "failed";

export interface PlannerOutput {
  summary: string;
  priority: "low" | "medium" | "high";
  steps: string[];
  risks: string[];
}

export interface FileChange {
  path: string;
  action: "create" | "modify" | "delete";
  content: string;
}

export interface CoderOutput {
  explanation: string;
  changes: FileChange[];
}

export interface ReviewComment {
  file: string;
  line: number;
  comment: string;
  level: "info" | "warning" | "error";
}

export interface ReviewerOutput {
  approved: boolean;
  summary: string;
  comments: ReviewComment[];
}

export interface Run {
  id: string;
  job_id: string;
  repo_url: string;
  pr_url: string;
  status: RunStatus;
  current_step: string;
  diff: string;
  plan: PlannerOutput | null;
  metadata: any;
  created_at: string;
  updated_at: string;
}

export interface RunsResponse {
  runs: Run[];
  limit: number;
  offset: number;
}

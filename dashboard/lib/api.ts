import { Run, RunsResponse } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function getRuns(limit = 10, offset = 0): Promise<RunsResponse> {
  const res = await fetch(`${API_BASE_URL}/runs?limit=${limit}&offset=${offset}`, {
    next: { revalidate: 10 },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch runs");
  }
  const data = await res.json();
  return data;
}

export async function getRun(id: string): Promise<Run> {
  console.log(id)
  const res = await fetch(`${API_BASE_URL}/runs/${id}`, {
    next: { revalidate: 5 },
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch run ${id}`);
  }
  const data = await res.json();
  return data;
}

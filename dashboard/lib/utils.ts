export function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function getStatusVariant(status: string) {
  switch (status) {
    case "completed":
      return "success";
    case "failed":
      return "error";
    case "pending":
      return "default";
    case "awaiting_approval":
      return "warning";
    default:
      return "info";
  }
}

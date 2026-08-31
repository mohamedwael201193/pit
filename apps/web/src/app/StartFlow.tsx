import { Navigate } from "react-router-dom";

/** Legacy /app/start. Pairing is step 1. */
export function StartFlow() {
  return <Navigate to="/pair" replace />;
}

import type { IFundReleaseActor } from "@/redux/api/trustees/interface";
import ActorAvatar from "./ActorAvatar";

export default function FundReleaseActor({
  actor,
  fallback = "N/A",
}: {
  actor?: IFundReleaseActor | null;
  fallback?: string;
}) {
  const name =
    actor?.name?.trim() ||
    [actor?.first_name, actor?.last_name].filter(Boolean).join(" ").trim();
  const email = actor?.email?.trim();
  const organization = actor?.organization_name?.trim();
  return (
    <div
      style={{ display: "flex", alignItems: "center", gap: 8, marginTop: 4 }}
    >
      <ActorAvatar name={name || email || organization} />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: 4,
          minWidth: 0,
          overflowWrap: "anywhere",
        }}
      >
        <span>{name || email || organization || fallback}</span>
        {organization && (name || email) && <small>{organization}</small>}
        {/* {email && name && <small>{email}</small>} */}
      </div>
    </div>
  );
}

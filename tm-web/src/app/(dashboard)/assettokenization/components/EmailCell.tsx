"use client";
import { useFetchUserByCriteriaQuery } from "@/redux/api/users";

interface EmailCellProps {
  initiatorUsername: string;
}

const EmailCell: React.FC<EmailCellProps> = ({ initiatorUsername }) => {
  // The hook is now safely inside a component, ensuring a consistent call order.
  const { data } = useFetchUserByCriteriaQuery(
    { username: initiatorUsername || "" },
    { skip: !initiatorUsername }
  );

  return <span>{data?.data?.user_info?.email || "Loading..."}</span>;
};

export default EmailCell;

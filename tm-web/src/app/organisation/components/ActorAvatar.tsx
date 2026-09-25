import { getInitials } from "@/utils/getInitials";
import styled from "styled-components";

export default function ActorAvatar({ name }: { name?: string | null }) {
  return <Circle aria-hidden="true">{getInitials(name)}</Circle>;
}

const Circle = styled.div`
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 10px;
  font-weight: 600;
  color: #191919;
  border-radius: 50%;
  background: #acd1ef;
`;

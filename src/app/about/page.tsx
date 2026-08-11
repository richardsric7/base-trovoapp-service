import { Metadata } from "next";
import AboutPage from "@/features/about";

export const metadata: Metadata = {
  title: "About Us | Trovotech",
  description:
    "Learn more about Trovotech's mission, vision, values, regulatory alignment, and the experienced team building infrastructure for real-world asset tokenization.",
};

export default function About() {
  return <AboutPage />;
}

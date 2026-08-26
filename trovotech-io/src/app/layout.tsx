import type { Metadata } from "next";
import { Montserrat } from "next/font/google";
import "./globals.css";
import Navbar from "@/component/Navbar";
import Footer from "@/component/Footer";

const montserrat = Montserrat({
  variable: "--font-montserrat",
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700", "800", "900"],
});

export const metadata: Metadata = {
  title: "Trovotech — Regulated Real-World Asset Tokenization",
  description:
    "Trovotech turns productive real-world assets into compliant digital tokens — fractional ownership, global investment, and borderless liquidity on a regulated blockchain network. Registered with the SEC of Nigeria.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${montserrat.variable} h-full antialiased`}>
      <head>
        <link rel="stylesheet" href="/font/stylesheet.css" />
      </head>
      <body className="min-h-full flex flex-col bg-white">
        <Navbar />
        {children}
        <Footer />
      </body>
    </html>
  );
}

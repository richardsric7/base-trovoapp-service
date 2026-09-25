import type { Metadata } from "next";
import { Montserrat } from "next/font/google";
import "./globals.css";
import StyledComponentsRegistry from "@/lib/registry";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { StoreProvider } from "@/redux";
import TokenRefresher from "./TokenRefresher";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import { ConfigProvider } from "antd";

// const inter = Inter({ subsets: ["latin"] });

const montserrat = Montserrat({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-montserrat",
});

export const metadata: Metadata = {
  title: "Trovotech Manager",
  description: "Trovo Manager",
  icons: {
    icon: "/TransperantLogo12.svg",
    shortcut: "/TransperantLogo12.svg",
    apple: "/TransperantLogo12.svg",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${montserrat.className} ${montserrat.variable}`}>
        <AntdRegistry>
          <ConfigProvider
            theme={{
              token: {
                fontFamily: "var(--font-montserrat), Montserrat, sans-serif",
              },
            }}
          >
            <TokenRefresher />
            <link rel="icon" href="/TransperantLogo12.svg" sizes="any" />
            <StyledComponentsRegistry>
              <StoreProvider>
                <ToastContainer />
                {children}
              </StoreProvider>
            </StyledComponentsRegistry>
          </ConfigProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}

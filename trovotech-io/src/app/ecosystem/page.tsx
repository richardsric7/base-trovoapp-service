"use client";

import React from "react";
import Link from "next/link";
import { FaApple, FaGooglePlay } from "react-icons/fa";
import PrimaryButton from "@/component/PrimaryButton";

const EcosystemPage = () => {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Header Section */}
      <section className="bg-light3 pt-24 sm:pt-32 pb-16 sm:pb-20 text-center px-5 sm:px-8">
        <div className="max-w-4xl mx-auto">
          <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 text-xs sm:text-sm font-matahari">
            THE TROVOTECH ECOSYSTEM
          </p>
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-extrabold leading-[1.1] mb-8 text-secondary font-matahari">
            The Infrastructure Layer for Digital Capital Markets
            {/* Every piece of the puzzle.
            <br className="hidden sm:block" /> One connected system. */}
          </h1>
          <p className="text-gray-600 text-base sm:text-lg max-w-3xl mx-auto leading-relaxed">
            Trovotech is a unified ecosystem of digital infrastructure
            purpose-built to modernize how productive assets are financed,
            administered, and exchanged. Through three tightly integrated
            platforms, we enable regulated capital formation, transparent
            lifecycle management, institutional governance, and compliant
            secondary market participation—creating a seamless experience for
            issuers, investors, and ecosystem partners.
          </p>
        </div>
      </section>

      {/* Main Content Sections */}
      <section className="py-20 px-5 sm:px-8 md:px-12 lg:px-20 max-w-[1400px] mx-auto space-y-24 sm:space-y-32 w-full">
        {/* Section 1: Trovo App */}
        <div className="flex flex-col md:flex-row items-center gap-10 lg:gap-24">
          <div className="md:w-1/2 flex flex-col items-start text-left">
            <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 text-xs sm:text-sm font-matahari">
              CAPITAL FORMATION PLATFORM
            </p>
            <h2 className="text-4xl md:text-5xl font-extrabold mb-6 text-secondary font-matahari">
              Trovo App
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed text-[15px] sm:text-base">
              Trovotech's primary capital formation platform, the digital
              infrastructure layer through which real-world productive assets
              are structured, issued, monitored, and transacted within a
              regulated environment. It features non custodial multi-signature
              wallets capability, multi-sub-wallets and a range of advanced
              feature for the secure storage and management of digital assets.
            </p>

            <PrimaryButton
              label="Explore the platform"
              href="https://trovo.app/"
              className="!text-white"
            />
            {/* <div className="flex flex-wrap gap-4">
              <a
                href="https://play.google.com/store/apps/details?id=com.trovo.wallet&pcampaignid=web_share"
                className="flex items-center gap-3 bg-primary hover:bg-primart/80 transition-colors text-white px-5 py-2.5 rounded-lg"
              >
                <FaGooglePlay className="text-2xl text-white" />
                <div className="text-left flex flex-col">
                  <span className="text-[10px] leading-tight uppercase font-medium text-white">
                    GET IT ON
                  </span>
                  <span className="font-semibold text-[15px] leading-tight text-white">
                    Google Play
                  </span>
                </div>
              </a>
              <a
                href="#"
                className="flex items-center gap-3 bg-[#111] hover:bg-black transition-colors text-white px-5 py-2.5 rounded-lg"
              >
                <FaApple className="text-2xl mb-1 text-white" />
                <div className="text-left flex flex-col">
                  <span className="text-[10px] leading-tight font-medium text-white">
                    Download on the
                  </span>
                  <span className="font-semibold text-[15px] leading-tight text-white">
                    App Store
                  </span>
                </div>
              </a>
            </div> */}
          </div>
          <div className="md:w-1/2 flex justify-center w-full">
            <img
              src="/img/new-website/trovoapp.png"
              alt="Trovo App Mockup"
              className="w-full h-auto object-contain max-h-[500px]"
            />
          </div>
        </div>

        {/* Section 2: TrovoP2P */}
        <div className="flex flex-col md:flex-row-reverse items-center gap-10 lg:gap-24">
          <div className="md:w-1/2 flex flex-col items-start text-left">
            <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 text-xs sm:text-sm font-matahari">
              SECONDARY MARKET INFRASTRUCTURE
            </p>
            <h2 className="text-4xl md:text-5xl font-extrabold mb-6 text-secondary font-matahari">
              TrovoP2P
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed text-[15px] sm:text-base">
              Trovotech's secondary transfer infrastructure, facilitating
              compliant peer-to-peer transfer of eligible digital assets. It
              supports structured liquidity, capital recycling, and transparent
              transfer within defined regulatory parameters. it is equipped with
              a platform managed escrow ensuring safe and secure transactions
            </p>
            <PrimaryButton
              label="Explore the platform"
              href="https://p2p.trovo.app/"
              className="!text-white"
            />
          </div>
          <div className="md:w-1/2 flex justify-center w-full">
            <img
              src="/img/new-website/trovop2p.png"
              alt="TrovoP2P Interface"
              className="w-full h-auto object-contain rounded-xl "
            />
          </div>
        </div>

        {/* Section 3: Trovo Manager */}
        <div className="flex flex-col md:flex-row items-center gap-10 lg:gap-24">
          <div className="md:w-1/2 flex flex-col items-start text-left">
            <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 text-xs sm:text-sm font-matahari">
              ASSET LIFECYCLE MANAGEMENT
            </p>
            <h2 className="text-4xl md:text-5xl font-extrabold mb-6 text-secondary font-matahari">
              Trovo Manager
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed text-[15px] sm:text-base">
              The internal Institutional-grade infrastructure that integrates
              established legal structures, governance arrangements, fiduciary
              oversight, regulatory compliance, professional intermediaries, and
              technology enabled infrastructure into a unified administrative
              dashboard overseeing the entire Trovotech ecosystem.
              <br />
              <br />
              It gives the trustees, custodians, asset managers and other
              organizations and stakeholders involved in the transactions
              controlled, auditable access to administer, and report on
              tokenized assets throughout their full lifecycle.
            </p>
          </div>
          <div className="md:w-1/2 flex justify-center w-full">
            {/* <div className="bg-light3 rounded-[40px] w-full max-w-[550px] flex items-center justify-center p-8 sm:p-12"> */}
            <img
              src="/img/new-website/trovomanager.png"
              alt="Trovo Manager Dashboard"
              className="w-full h-auto object-contain"
            />
            {/* </div> */}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-dark text-white pt-24 pb-8 px-5 text-center mt-auto">
        <div className="max-w-3xl mx-auto border-b border-white/10 pb-16">
          <h2 className="text-3xl sm:text-4xl md:text-[40px] font-bold mb-6 font-matahari">
            Start Building on Trovotech
          </h2>
          <p className="text-gray-400 mb-10 text-sm sm:text-base">
            Leave the details to integrated APIs or discuss non-technical
            deployments <br className="hidden sm:block" /> for your institution.
          </p>
          <div className="flex flex-wrap justify-center gap-4 mb-8">
            <Link
              href="https://calendly.com/thomasenechi"
              className="px-8 py-3 bg-primary hover:bg-primary-dark transition-colors rounded-lg font-medium"
              target="_blank"
            >
              Tokenize an asset
            </Link>
            <Link
              href="https://p2p.trovo.app/"
              className="px-8 py-3 bg-transparent border border-white/30 hover:border-white transition-colors rounded-lg font-medium"
              target="_blank"
            >
              Trade an asset
            </Link>
          </div>
          <p className="text-gray-400 text-sm">
            Speak with our team?{" "}
            <Link
              href="/contact"
              className="text-white hover:text-primary transition-colors border-b border-white hover:border-primary pb-0.5"
            >
              Contact us
            </Link>
          </p>
        </div>
      </section>
    </div>
  );
};

export default EcosystemPage;

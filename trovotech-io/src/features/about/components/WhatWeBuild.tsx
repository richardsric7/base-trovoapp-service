import React from "react";
import Image from "next/image";
import abtImg from "../../../../public/img/new-website/abtimg.jpg";
import { FaCheckCircle } from "react-icons/fa";
import { FaAward } from "react-icons/fa6";
import { RiVerifiedBadgeFill } from "react-icons/ri";

const PARTICIPANTS = [
  "TRUSTEES",
  "CUSTODIANS",
  "ASSET MANAGERS",
  "ISSUING HOUSES",
];

const BENEFITS = [
  {
    title: "Regulated Infrastructure",
    desc: "Built within compliant capital market frameworks.",
  },
  {
    title: "Transparent Investment Access",
    desc: "Improved visibility into asset structures and performance.",
  },
  {
    title: "Broader Participation",
    desc: "Enabling wider access to productive investment opportunities.",
  },
  {
    title: "Institutional Governance",
    desc: "Integrated trustees, custodians, and asset managers.",
  },
];



const WhatWeBuild = () => {
  return (
    <section className="bg-white">
      {/* 1. Header & Intro */}
      <div className="max-w-[1280px] mx-auto px-5 sm:px-8 md:px-12 lg:px-20 pt-16 md:pt-24 pb-12">
        <div className="text-center">
          <p className="text-[#007cdf] font-bold tracking-[0.15em] uppercase text-xs sm:text-sm font-matahari mb-4">
            WHAT WE BUILD
          </p>
          <h2 className="font-matahari text-3xl sm:text-4xl md:text-5xl lg:text-[44px] font-extrabold text-secondary leading-[1.2] tracking-tight mb-6 max-w-[850px] mx-auto">
            Infrastructure for the Tokenized Asset Economy
          </h2>
          <p className="text-gray-600 text-sm sm:text-base leading-[1.65] max-w-[760px] mx-auto font-sans font-normal">
            Trovotech provides regulated infrastructure that enables real-world
            assets such as infrastructure projects, real estate, agriculture,
            investment funds, etc. to be structured into compliant digital
            financial instruments. The platform integrates key capital market
            participants to support transparent and governed investment
            participation. They include:
          </p>
        </div>
      </div>

      {/* 2. Capital Market Participants Bar */}
      <div className="w-full bg-[#111e35] py-6 px-5 sm:px-8">
        {/* Mobile Grid */}
        <div className="md:hidden grid grid-cols-2 gap-y-4 gap-x-2 max-w-md mx-auto justify-items-center">
          {PARTICIPANTS.map((item, idx) => (
            <div key={idx} className="flex items-center gap-2.5">
              <RiVerifiedBadgeFill color="#00A859" />
              <span className="text-white text-[11px] sm:text-xs font-extrabold tracking-wider uppercase font-sans">
                {item}
              </span>
            </div>
          ))}
        </div>

        {/* Desktop Flex Bar */}
        <div className="hidden md:flex justify-center items-center">
          {PARTICIPANTS.map((item, idx) => (
            <div
              key={idx}
              className="flex items-center gap-3 px-6 lg:px-10 border-r border-slate-700/50 last:border-none"
            >
              <RiVerifiedBadgeFill color="#00A859" />
              <span className="text-white text-xs lg:text-sm font-extrabold tracking-widest uppercase font-sans">
                {item}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* 3. Key Benefits Section */}
      <div className=" mx-auto px-5 sm:px-8 md:px-12 lg:px-20 pt-16 md:pt-24">
        <div className="text-center mb-10">
          <p className="text-[#007cdf] font-bold tracking-[0.15em] uppercase text-xs sm:text-sm font-matahari">
            WHAT ARE THE KEY BENEFITS?
          </p>
        </div>

        {/* Benefits Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 sm:gap-8 mb-20 md:mb-28">
          {BENEFITS.map((benefit, idx) => (
            <div
              key={idx}
              className="flex flex-col justify-between p-8 rounded-3xl bg-[#f2f6f9] hover:bg-[#ebf0f5] transition-all duration-300 hover:shadow-lg hover:-translate-y-1 group min-h-[220px] sm:min-h-[240px]"
            >
              <h3 className="text-lg sm:text-xl font-extrabold text-secondary font-matahari leading-snug tracking-tight mb-8">
                {benefit.title}
              </h3>
              <p className="text-gray-600 text-xs sm:text-sm leading-relaxed mt-auto font-sans font-normal">
                {benefit.desc}
              </p>
            </div>
          ))}
        </div>

        {/* 4. Why This Matters Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-16 items-center pb-20 md:pb-28">
          {/* Left Text */}
          <div className="flex flex-col justify-center">
            <p className="text-[#007cdf] font-bold tracking-[0.15em] uppercase text-xs sm:text-sm font-matahari mb-4">
              WHY THIS MATTERS
            </p>
            <h2 className="text-3xl sm:text-4xl md:text-5xl lg:text-[44px] font-extrabold text-secondary font-matahari leading-[1.2] tracking-tight mb-6">
              Expanding Access to Wealth Creation
            </h2>
            <div className="space-y-5 text-gray-600 text-sm sm:text-base leading-[1.65] font-sans font-normal">
              <p>
                For decades, access to institutional-quality investment
                opportunities has remained limited to a small segment of the
                population.
              </p>
              <p>
                Trovotech&apos;s long-term vision is to help create a financial
                system where productive assets become more accessible,
                transparent, and investable across emerging markets. We believe
                capital markets should not only serve institutions, but also
                enable broader participation in economic growth.
              </p>
            </div>
          </div>

          {/* Right Image */}
          <div className="relative w-full aspect-[4/3] sm:aspect-[1.1] md:aspect-[1.2] lg:aspect-[1.1] xl:aspect-[1.15] rounded-[32px] overflow-hidden shadow-md">
            <Image
              src={abtImg}
              alt="Expanding access to wealth creation"
              fill
              className="object-cover"
              priority
              sizes="(max-width: 768px) 100vw, 50vw"
            />
          </div>
        </div>
      </div>
    </section>
  );
};

export default WhatWeBuild;

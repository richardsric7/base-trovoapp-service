"use client";

import { useState } from "react";
import { IoClose } from "react-icons/io5";

const COMPLIANCE_CARDS = [
  {
    icon: "/img/new-website/sec-logo.svg",
    title:
      "Regulated by the Securities and Exchange Commission of Nigeria (SEC)",
    short:
      "Trovotech is regulated under the Real World Asset Tokenization and Offering Platform category...",
    full: "Trovotech is regulated under the Real World Asset Tokenization and Offering Platform (RATOP) category of the Securities and Exchange Commission of Nigeria. This means the platform operates within a defined legal and supervisory framework — not in a grey area. Every investment offering must meet regulatory standards before it reaches investors.",
  },
  {
    icon: "/img/new-website/verified.svg",
    title: "Every participant is verified",
    short:
      "Before anyone can invest or transact on Trovotech, they must complete identity verification. This...",
    full: "Before anyone can invest or transact on Trovotech, they must complete identity verification. This applies to all users — not just institutions. Anti-Money Laundering and Counter-Terrorism Financing controls are embedded into the platform, with all transactions monitored. This protects the entire ecosystem, not just individual users.",
  },
  {
    icon: "/img/new-website/Group.svg",
    title: "Trovotech never holds your money",
    short:
      "Investor funds are held by licensed, independent custodians — not by Trovotech. Trovotech provides the...",
    full: "Investor funds are held by licensed, independent custodians — not by Trovotech. Trovotech provides the infrastructure that coordinates investment; the actual safekeeping of capital belongs to regulated institutions with legal obligations to protect it. Even if the platform were to stop operating, your funds would remain with the custodian.",
  },
  {
    icon: "/img/new-website/law.svg",
    title: "Your rights exist in law, not on ledger",
    short:
      "Every asset on Trovotech is legally wrapped in a trust or Special Purpose Vehicle before it reaches i...",
    full: "Every asset on Trovotech is legally wrapped in a trust or Special Purpose Vehicle before it reaches investors. This means your economic rights — income, participation, or repayment — are defined in legal documents and enforced by independent trustees. If Trovotech disappeared tomorrow, those rights would still exist and be enforceable.",
  },
];

const Compliance = () => {
  const [expandedIndex, setExpandedIndex] = useState<number | null>(0); // First card is active by default in index.html!

  const handleToggle = (index: number) => {
    if (expandedIndex === index) {
      setExpandedIndex(null);
    } else {
      setExpandedIndex(index);
    }
  };

  return (
    <section className="py-20  bg-light">
      <div className="page-container">
        <div className="   mb-12">
          <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 font-matahari text-sm sm:text-base">
            REGULATION & TRUST
          </p>
          <h2 className="text-4xl sm:text-5xl 3xl:text-[56px] font-extrabold leading-[1.1] mb-6 text-secondary font-matahari">
            Built on compliance, <br />
            not assumptions
          </h2>
          <p className="text-gray-600 leading-[160%] text-sm sm:text-base 3xl:text-lg font-normal w-[70%] 3xl:w-[70%]">
            Most digital asset platforms ask you to trust the technology.
            Trovotech operates within a legal and regulatory framework — so that
            trust is backed by enforceable structures, licensed institutions,
            and government oversight.
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 3xl:gap-10 mx-auto w-full">
          {COMPLIANCE_CARDS.map((card, idx) => {
            const isExpanded = expandedIndex === idx;
            return (
              <div
                key={idx}
                className={`relative flex flex-col p-6 3xl:p-7 rounded-3xl transition-all duration-300 min-h-[320px] 3xl:min-h-[350px] w-full border ${
                  isExpanded
                    ? "bg-secondary border-secondary shadow-xl scale-[1.02]"
                    : "bg-[#f2f6f9] border-[#e2e8f0] hover:shadow-lg hover:-translate-y-1"
                }`}
              >
                {/* Icon wrapper */}
                <div
                  className={`w-14 h-14 rounded-full flex items-center justify-center mb-6 transition-colors duration-300 ${
                    isExpanded
                      ? "bg-white text-secondary"
                      : "bg-white text-primary"
                  }`}
                >
                  <img
                    src={card.icon}
                    alt={card.title}
                    className="w-7 h-7 object-cover"
                  />
                </div>

                {/* Content */}
                <h3
                  className={`text-lg sm:text-xl font-bold mb-4 font-matahari tracking-tight line-clamp-3 mt-auto ${
                    isExpanded ? "text-white" : "text-secondary"
                  }`}
                >
                  {card.title}
                </h3>

                <p
                  className={`text-sm leading-[150%] mb-6 ${
                    isExpanded ? "text-white/90" : "text-trovo-black"
                  }`}
                >
                  {isExpanded ? card.full : card.short}
                </p>

                {/* Toggle button */}
                <button
                  onClick={() => handleToggle(idx)}
                  className={`w-fit border-none rounded-full px-3.5 py-1.5 text-xs font-semibold cursor-pointer transition-all duration-300 ${
                    isExpanded
                      ? "bg-white text-secondary hover:bg-white/90"
                      : "bg-[#acd1ef] text-trovo-black hover:bg-[#d8ebff]"
                  }`}
                  type="button"
                >
                  {isExpanded ? (
                    <span className="flex items-center gap-1">
                      <IoClose />
                      Close
                    </span>
                  ) : (
                    "Learn More"
                  )}
                </button>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
};

export default Compliance;

"use client";

import React, { useState, useEffect, useRef } from "react";

interface Feature {
  title: string;
  desc: string;
  image: string;
  alt: string;
  icon: string;
}

interface Tab {
  id: string;
  label: string;
  features: Feature[];
}

const MARKET_TABS: Tab[] = [
  {
    id: "issuers",
    label: "Investors",
    features: [
      {
        title: "Institutional Investors",
        desc: "Access regulated structured financial instruments across real estate, infrastructure, agriculture, and private funds — previously inaccessible at scale.",
        image: "/img/new-website/img1.jpg",
        alt: "Access to global capital",
        icon: "/img/new-website/briefcase.svg",
      },
      {
        title: "High Net Worth Individuals",
        desc: "Invest in income-generating real assets with defined legal rights, transparent reporting, and secondary market access via TrovoP2P.",
        image: "/img/new-website/img2.png",
        alt: "Lower cost of issuance",
        icon: "/img/new-website/crown.svg",
      },
      {
        title: "Retail Participants",
        desc: "Invest in income-generating real assets with defined legal rights, transparent reporting, and secondary market access via TrovoP2P.",
        image: "/img/new-website/img3.jpg",
        alt: "Instant settlement",
        icon: "/img/new-website/user.svg",
      },
    ],
  },
  {
    id: "institutional",
    label: "Asset Issuers ",
    features: [
      {
        title: "Real Estate Developers",
        desc: "Raise structured capital against your real estate portfolio through compliant issuance — without replacing your existing legal arrangements.",
        image: "/img/new-website/img4.png",
        alt: "Custodial solutions",
        icon: "/img/new-website/house.svg",
      },
      {
        title: "Infrastructure Sponsors",
        desc: "Finance energy, transport, and utility projects through long-term structured instruments distributed to a broad investor base.",
        image: "/img/new-website/img5.png",
        alt: "Reporting and analytics",
        icon: "/img/new-website/fa6-solid_bridge.svg",
      },
      {
        title: "Agricultural Businesses",
        desc: "Access patient capital for crop production, processing, and distribution through seasonal structures that match your cash flow cycles.",
        image: "/img/new-website/img6.png",
        alt: "Reporting and analytics",
        icon: "/img/new-website/flower.svg",
      },
    ],
  },
  {
    id: "individual",
    label: "Institutions",
    features: [
      {
        title: "Trustees and Custodians",
        desc: "Integrate with Trovotech's infrastructure to serve as legal stewards of investment structures — protecting investor rights end-to-end.",
        image: "/img/new-website/img4.png",
        alt: "Fractional ownership",
        icon: "/img/new-website/shield.svg",
      },
      {
        title: "Asset Managers",
        desc: "Use Trovo Manager to administer fund lifecycles, reporting obligations, and investor communications from a single governed platform.",
        image: "/img/new-website/img5.png",
        alt: "24/7 market",
        icon: "/img/new-website/uis_chart.svg",
      },
      {
        title: "Banks and Issuing Houses",
        desc: "Distribute structured investment products to your client base, provide custody services, and participate in asset financing — all via API.",
        image: "/img/new-website/img4.png",
        alt: "24/7 market",
        icon: "/img/new-website/streamline_bank-solid.svg",
      },
    ],
  },
  {
    id: "regulators",
    label: "Partners",
    features: [
      {
        title: "On-chain transparency",
        desc: "Use Trovo Manager to administer fund lifecycles, reporting obligations, and investor communications from a single governed platform.",
        image: "/img/new-website/img4.png",
        alt: "On-chain transparency",
        icon: "/img/new-website/user.svg",
      },
      {
        title: "Automated reporting",
        desc: "Distribute structured investment products to your client base, provide custody services, and participate in asset financing — all via API.",
        image: "/img/new-website/img5.png",
        alt: "Automated reporting",
        icon: "/img/new-website/user.svg",
      },
      {
        title: "Automated reporting",
        desc: "Distribute structured investment products to your client base, provide custody services, and participate in asset financing — all via API.",
        image: "/img/new-website/img5.png",
        alt: "Automated reporting",
        icon: "/img/new-website/uis_chart.svg",
      },
    ],
  },
];

const FEATURE_DURATION = 2400; // ms

const Market = () => {
  const [activeTabIdx, setActiveTabIdx] = useState(0);
  const [activeFeatureIdx, setActiveFeatureIdx] = useState(0);
  const [progress, setProgress] = useState(0);

  const containerRef = useRef<HTMLDivElement | null>(null);
  const firstFeatureRef = useRef<HTMLDivElement | null>(null);
  const lastFeatureRef = useRef<HTMLDivElement | null>(null);
  const [imageStyle, setImageStyle] = useState<React.CSSProperties>({});

  const tab = MARKET_TABS[activeTabIdx];
  const activeFeature = tab.features[activeFeatureIdx];

  // Dynamic Image Alignment calculations
  const calculateLayout = () => {
    if (window.innerWidth <= 768) {
      setImageStyle({});
      return;
    }

    if (
      containerRef.current &&
      firstFeatureRef.current &&
      lastFeatureRef.current
    ) {
      const containerRect = containerRef.current.getBoundingClientRect();
      const firstRect = firstFeatureRef.current.getBoundingClientRect();
      const lastRect = lastFeatureRef.current.getBoundingClientRect();

      const topOffset = Math.max(0, firstRect.top - containerRect.top);
      const stackHeight = Math.max(0, lastRect.bottom - firstRect.top);

      setImageStyle({
        marginTop: `${topOffset}px`,
        height: `${stackHeight}px`,
      });
    }
  };

  useEffect(() => {
    calculateLayout();
    window.addEventListener("resize", calculateLayout);
    return () => window.removeEventListener("resize", calculateLayout);
  }, [activeTabIdx]);

  // Autoplay Loop logic using requestAnimationFrame
  useEffect(() => {
    let animId: number;
    let startTime = performance.now();

    const loop = (now: number) => {
      const elapsed = now - startTime;
      const computedProgress = (elapsed / FEATURE_DURATION) * 100;

      if (elapsed >= FEATURE_DURATION) {
        setProgress(0);
        const nextFeatureIdx = activeFeatureIdx + 1;
        if (nextFeatureIdx < tab.features.length) {
          setActiveFeatureIdx(nextFeatureIdx);
          startTime = performance.now();
          animId = requestAnimationFrame(loop);
        } else {
          // Switch to next tab
          const nextTabIdx = (activeTabIdx + 1) % MARKET_TABS.length;
          setActiveTabIdx(nextTabIdx);
          setActiveFeatureIdx(0);
          setProgress(0);
        }
      } else {
        setProgress(computedProgress);
        animId = requestAnimationFrame(loop);
      }
    };

    animId = requestAnimationFrame(loop);
    return () => cancelAnimationFrame(animId);
  }, [activeTabIdx, activeFeatureIdx]);

  const handleTabClick = (idx: number) => {
    setActiveTabIdx(idx);
    setActiveFeatureIdx(0);
    setProgress(0);
  };

  return (
    <section className="py-20 bg-light2">
      <div className="page-container">
        <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 font-matahari text-sm sm:text-base">
          WHO IS THIS FOR
        </p>
        <div className="max-w-[700px] 3xl:max-w-[780px] mb-12">
          <h2 className="text-4xl sm:text-5xl 3xl:text-[56px] font-extrabold leading-[1.1] mb-6 text-secondary font-matahari">
            Built for every side of <br />
            the market
          </h2>
        </div>

        {/* Tab Buttons */}
        <div className="flex flex-wrap gap-4 mb-16">
          {MARKET_TABS.map((t, idx) => (
            <button
              key={t.id}
              onClick={() => handleTabClick(idx)}
              className={`inline-flex items-center gap-2 px-6 py-3 rounded-lg font-medium cursor-pointer transition-all duration-300 ${
                activeTabIdx === idx
                  ? "bg-secondary text-white shadow-md scale-105"
                  : "bg-trovo-blue text-secondary hover:bg-gray-200"
              }`}
              type="button"
            >
              <span>{t.label}</span>
            </button>
          ))}
        </div>

        {/* Active Content Grid */}
        <div
          ref={containerRef}
          className="grid grid-cols-1 md:grid-cols-2 gap-8 lg:gap-16 3xl:gap-20 items-start bg-white rounded-2xl p-6 3xl:p-8 shadow-sm border border-gray-100"
        >
          <div className="market-info">
            <div className="market-features space-y-4">
              {tab.features.map((feature, fIdx) => {
                const isActive = activeFeatureIdx === fIdx;
                return (
                  <div
                    key={fIdx}
                    ref={(el) => {
                      if (fIdx === 0) firstFeatureRef.current = el;
                      if (fIdx === tab.features.length - 1)
                        lastFeatureRef.current = el;
                    }}
                    className={`relative overflow-hidden flex flex-col gap-3 p-4 3xl:p-5 rounded-xl transition-all duration-300 border ${
                      isActive
                        ? "bg-[#1f3255] border-[#1f3255] text-white translate-x-1"
                        : "bg-[#e8ecf4] border-transparent text-[#22344f] hover:border-gray-300 hover:translate-x-0.5"
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-full bg-white flex items-center justify-center flex-shrink-0">
                        <img
                          src={feature.icon}
                          alt={feature.title}
                          className="w-6 h-6 object-contain"
                        />
                      </div>
                      <span
                        className={`text-lg font-bold leading-tight font-sans transition-colors ${
                          isActive ? "text-white" : "text-[#22344f]"
                        }`}
                      >
                        {feature.title}
                      </span>
                    </div>

                    <p
                      className={`text-sm leading-[150%] transition-all duration-500 overflow-hidden font-normal ${
                        isActive
                          ? "max-h-[200px] opacity-100 mt-2 text-white/88"
                          : "max-h-0 opacity-0 pointer-events-none"
                      }`}
                    >
                      {feature.desc}
                    </p>

                    {/* Progress Bar inside card */}
                    <div
                      className={`absolute left-2 right-2 bottom-1.5 h-0.5 bg-white rounded-full origin-left transition-transform duration-75 ${
                        isActive ? "opacity-100" : "opacity-0"
                      }`}
                      style={{
                        transform: `scaleX(${isActive ? progress / 100 : 0})`,
                      }}
                    />
                  </div>
                );
              })}
            </div>
          </div>

          <div
            className="w-full relative rounded-2xl overflow-hidden shadow-md flex items-center"
            style={imageStyle}
          >
            <img
              src={activeFeature.image}
              alt={activeFeature.alt}
              className="w-full h-full object-cover rounded-2xl transition-all duration-500 max-h-[300px] md:max-h-none"
            />
          </div>
        </div>
      </div>
    </section>
  );
};

export default Market;

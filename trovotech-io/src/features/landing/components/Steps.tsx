"use client";

import { useState, useEffect, useRef } from "react";

const STEPS_DATA = [
  {
    number: 1,
    title: "Asset Structuring",
    description:
      "The asset sponsor submits an asset. Legal and financial teams structure it into a compliant financial instrument via a trust or SPV.",
  },
  {
    number: 2,
    title: "Validation and Admission",
    description:
      "Assets undergo legal validation of ownership, financial assessment, risk identification, and regulatory admission review.",
  },
  {
    number: 3,
    title: "Digital Issuance",
    description:
      "Compliant digital tokens are issued representing defined financial interests — income rights, ownership, or repayment obligations.",
  },
  {
    number: 4,
    title: "Lifecycle and liquidity",
    description:
      "Investors access, hold, and trade their positions. Trustees and custodians ensure ongoing compliance and fund protection.",
  },
];

const Steps = () => {
  const [activeIndex, setActiveIndex] = useState(0);
  const sectionRef = useRef<HTMLElement>(null);
  const rafRef = useRef<number | null>(null);
  const smoothProgressRef = useRef(0);
  const rawProgressRef = useRef(0);

  useEffect(() => {
    const handleScroll = () => {
      if (!sectionRef.current) return;
      const rect = sectionRef.current.getBoundingClientRect();
      const sectionHeight = sectionRef.current.offsetHeight;
      const viewportHeight = window.innerHeight;
      const scrollableDistance = sectionHeight - viewportHeight;

      if (scrollableDistance <= 0) return;

      const scrolledIn = -rect.top;
      rawProgressRef.current = Math.max(
        0,
        Math.min(1, scrolledIn / scrollableDistance),
      );
    };

    const animate = () => {
      // Lerp smoothed progress toward the raw scroll progress (0.1 = ease factor)
      smoothProgressRef.current +=
        (rawProgressRef.current - smoothProgressRef.current) * 0.1;

      const newIndex = Math.min(
        STEPS_DATA.length - 1,
        Math.floor(smoothProgressRef.current * STEPS_DATA.length),
      );
      setActiveIndex(newIndex);

      rafRef.current = requestAnimationFrame(animate);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    window.addEventListener("resize", handleScroll);
    handleScroll();
    rafRef.current = requestAnimationFrame(animate);

    return () => {
      window.removeEventListener("scroll", handleScroll);
      window.removeEventListener("resize", handleScroll);
      if (rafRef.current !== null) cancelAnimationFrame(rafRef.current);
    };
  }, []);

  return (
    <section
      ref={sectionRef}
      className="relative bg-white py-0"
      style={{ height: `${STEPS_DATA.length * 25}vh` }}
    >
      {/* Sticky panel — fixed height, image cannot shake */}
      <div className="sticky top-0 h-screen flex items-center overflow-hidden">
        <div className="page-container py-0">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-24 3xl:gap-32 items-center">
            {/* Left: steps content */}
            <div>
              <p className="text-primary font-bold tracking-[0.1em] uppercase pt-4 mb-4 font-matahari text-sm sm:text-base">
                HOW IT WORKS
              </p>
              <h2 className="text-4xl sm:text-5xl 3xl:text-[56px] font-extrabold leading-[1.08] mb-6 text-secondary font-matahari">
                Four steps from asset to investor
              </h2>
              <p className="text-gray-600 mb-10 leading-[160%] text-sm sm:text-base 3xl:text-lg max-w-[520px] 3xl:max-w-[580px]">
                Trovotech provides the structured rails that take a real-world
                productive asset and make it legally investable at scale.
              </p>

              <div>
                {STEPS_DATA.map((step, idx) => (
                  <div
                    key={idx}
                    className="flex gap-4 sm:gap-6 px-4 py-4 3xl:py-5"
                  >
                    <div
                      className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 font-semibold text-sm transition-all duration-500 ${
                        activeIndex === idx
                          ? "bg-primary text-white scale-110 shadow-md"
                          : "bg-trovo-blue text-[#8f979d]"
                      }`}
                    >
                      {step.number}
                    </div>
                    <div className="flex-1">
                      <h3
                        className={`text-lg sm:text-xl font-bold transition-colors duration-500 font-sans ${
                          activeIndex === idx ? "text-secondary" : "text-grey"
                        }`}
                      >
                        {step.title}
                      </h3>
                      <p
                        className={`text-sm sm:text-base text-gray-600 leading-[150%] transition-all duration-700 ease-in-out origin-top overflow-hidden ${
                          activeIndex === idx
                            ? "max-h-[200px] opacity-100 mt-2"
                            : "max-h-0 opacity-0 pointer-events-none"
                        }`}
                      >
                        {step.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Right: image — inside fixed-height container, zero movement */}
            <div className="w-full h-[320px] sm:h-[420px] lg:h-[500px] 3xl:h-[560px] rounded-3xl overflow-hidden">
              <img
                src="/img/new-website/howitworks.jpg"
                alt="Investor using Trovo App"
                className="w-full h-full object-cover rounded-3xl"
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Steps;

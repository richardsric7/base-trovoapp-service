import Image from "next/image";
import image from "../../../../public/img/new-website/dotted-map.svg";
import { FaCheckCircle } from "react-icons/fa";

const VISION_CHECKLIST = [
  "Full suite of tokenized real-world assets",
  "Institutional-grade security",
  "Seamless API integration",
  "Multi-chain capability for maximum flexibility and reach",
];

const Vision = () => {
  return (
    <section className="py-20 md:py-28 bg-[#1D212F] text-white relative overflow-hidden">
      {/* Background Dotted Map */}
      <div className="absolute inset-y-0 left-0 -right-[60%] z-0 pointer-events-none overflow-hidden hidden md:block">
        <div className="max-w-[500px]  mx-auto  h-full flex items-start justify-end pt-20 md:pt-28">
          <Image
            src={image}
            alt="trovotech-vision-map"
            className="w-full max-w-[500px] md:max-w-[650px] h-auto object-contain object-right "
            priority
          />
        </div>
      </div>

      <div className="relative z-10  mx-auto px-5 sm:px-8 md:px-12 lg:px-20">
        {/* Header */}
        <div className="max-w-[850px] mb-12">
          <p className="text-[#7eb1e5] font-bold tracking-[0.15em] uppercase mb-4 text-xs sm:text-sm font-matahari">
            BEYOND THE PLATFORM
          </p>
          <h2 className="text-3xl sm:text-4xl md:text-5xl lg:text-[50px] font-extrabold leading-[1.15] mb-6 text-white font-matahari tracking-tight">
            Trovotech&apos;s vision extends <br />
            beyond a single platform.
          </h2>
        </div>

        {/* Green-checked capabilities list */}
        <div className="grid grid-cols-1 md:grid-cols-1 gap-x-12 gap-y-4 max-w-[850px] mb-16 font-sans">
          {VISION_CHECKLIST.map((item, index) => (
            <div key={index} className="flex items-start gap-3">
              <FaCheckCircle color="#00A859" />
              <span className="text-white/90 text-sm sm:text-base font-medium">
                {item}
              </span>
            </div>
          ))}
        </div>

        <div className="w-full h-px bg-white/10"></div>

        {/* Our Vision, Impact, & Mission Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-0 relative mt-12">
          {/* Vertical lines between columns (visible on md and up) */}
          <div className="hidden md:block absolute left-1/3 top-[-48px] bottom-0 w-px bg-white/10 -translate-x-1/2"></div>
          <div className="hidden md:block absolute left-2/3 top-[-48px] bottom-0 w-px bg-white/10 -translate-x-1/2"></div>

          {/* Card 1: Our Vision */}
          <div className="flex flex-col p-6 sm:p-8 md:pr-8 lg:pr-12">
            <h3 className="text-2xl font-extrabold text-white font-matahari mb-4 group-hover:text-[#7eb1e5] transition-colors duration-300">
              Our Vision
            </h3>
            <p className="text-white/70 text-sm sm:text-base leading-[175%] font-normal">
              To become the trusted capital infrastructure layer enabling
              structured investment into productive real-world assets across
              emerging markets.
              {/* To create a seamless, borderless financial ecosystem where real-world
              assets are digitized, democratized, and globally accessible, fostering
              economic growth and financial inclusion across emerging markets. */}
            </p>
          </div>

          {/* Horizontal line between rows (visible on mobile) */}
          <div className="h-px w-full bg-white/10 my-4 md:hidden"></div>

          {/* Card 2: Our Impact */}
          <div className="flex flex-col p-6 sm:p-8 md:px-8 lg:px-12">
            <h3 className="text-2xl font-extrabold text-white font-matahari mb-4 group-hover:text-[#7eb1e5] transition-colors duration-300">
              Our Impact
            </h3>
            <p className="text-white/70 text-sm sm:text-base leading-[175%] font-normal">
              To empower institutions, asset sponsors, and everyday investors by
              providing compliant, reliable, and secure tokenization rails that
              redefine wealth creation and asset management for the digital era.
            </p>
          </div>

          {/* Horizontal line between rows (visible on mobile) */}
          <div className="h-px w-full bg-white/10 my-4 md:hidden"></div>

          {/* Card 3: Our Mission */}
          <div className="flex flex-col p-6 sm:p-8 md:pl-8 lg:pl-12">
            <h3 className="text-2xl font-extrabold text-white font-matahari mb-4 group-hover:text-[#7eb1e5] transition-colors duration-300">
              Our Mission
            </h3>
            <p className="text-white/70 text-sm sm:text-base leading-[175%] font-normal">
              To engineer regulated, governance-driven digital market infrastructure that enables compliant issuance, lifecycle administration, and secondary liquidity for trustee-structured real-world assets while expanding disciplined access to capital and protecting market integrity.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Vision;

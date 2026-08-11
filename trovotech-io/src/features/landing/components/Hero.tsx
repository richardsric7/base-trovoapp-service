"use client";

import Image from "next/image";
import { useState, useEffect, useRef } from "react";

const HERO_IMAGES = [
  "/img/new-website/hero1.jpg",
  "/img/new-website/hero2.jpg",
  "/img/new-website/hero3.jpg",
];

const Hero = () => {
  const [lockedIndex, setLockedIndex] = useState(0);
  const [hoverIndex, setHoverIndex] = useState<number | null>(null);

  // Two-layer crossfade state
  const [bgA, setBgA] = useState(HERO_IMAGES[0]);
  const [bgB, setBgB] = useState(HERO_IMAGES[1]); // pre-fill so layer B is never empty
  const [activeLayer, setActiveLayer] = useState<"A" | "B">("A");

  const timerRef = useRef<NodeJS.Timeout | null>(null);

  // Target image we want to display
  const displayIndex = hoverIndex !== null ? hoverIndex : lockedIndex;
  const targetImage = HERO_IMAGES[displayIndex];

  // Triggers the crossfade swap when target changes
  useEffect(() => {
    const currentActiveImage = activeLayer === "A" ? bgA : bgB;
    if (targetImage !== currentActiveImage) {
      if (activeLayer === "A") {
        setBgB(targetImage);
        setActiveLayer("B");
      } else {
        setBgA(targetImage);
        setActiveLayer("A");
      }
    }
  }, [targetImage, activeLayer, bgA, bgB]);

  // Autoplay loop
  const startAutoplay = () => {
    stopAutoplay();
    timerRef.current = setInterval(() => {
      setLockedIndex((prev) => (prev + 1) % HERO_IMAGES.length);
    }, 5000);
  };

  const stopAutoplay = () => {
    if (timerRef.current) clearInterval(timerRef.current);
  };

  useEffect(() => {
    startAutoplay();
    return () => stopAutoplay();
  }, []);

  const handleCardClick = (index: number) => {
    setLockedIndex(index);
    startAutoplay();
  };

  return (
    <section className="relative md:h-screen overflow-hidden flex flex-col">
      {/* Background Crossfade Layers — Next.js <Image fill> so first image is preloaded */}
      <div className="absolute inset-0 z-1 ">
        {/* Layer A */}
        <div
          className={`absolute inset-0 transition-all duration-[1200ms] ease-in-out ${
            activeLayer === "A"
              ? "opacity-100 scale-100"
              : "opacity-0 scale-[1.04]"
          }`}
        >
          <Image
            src={bgA}
            alt=""
            fill
            sizes="100vw"
            className="object-cover object-top"
            // priority on the initial hero image so Next.js preloads it before paint
            priority={bgA === HERO_IMAGES[0]}
          />
        </div>

        {/* Layer B */}
        <div
          className={`absolute inset-0 transition-all duration-[1200ms] ease-in-out  ${
            activeLayer === "B"
              ? "opacity-100 scale-100"
              : "opacity-0 scale-[1.04]"
          }`}
        >
          <Image
            src={bgB}
            alt=""
            fill
            sizes="100vw"
            className="object-cover object-top"
          />
        </div>
      </div>

      {/* Layer Overlay Gradient */}
      <div className="absolute inset-0 z-2 hero-overlay-gradient" />

      {/* Hero Content Grid */}
      <div className="page-container relative z-5 flex-1 flex flex-col lg:flex-row items-center lg:items-center lg:justify-between pt-12 md:py-8">
        {/* <div className="page-container relative z-5 pt-28 3xl:pt-32  flex flex-col md:flex-row md:items-start md:justify-between flex-1"> */}
        <div className="w-full md:w-[560px] 3xl:w-[620px] text-white mt-8 inline-flex flex-col gap-4">
          <h1 className="text-white text-[36px]  md:text-[54px] xl:text-[64px] 3xl:text-[72px] font-extrabold leading-[1.1] md:leading-[72px] 3xl:leading-[80px] tracking-[0.17px] ">
            Stop watching <br />
            wealth get built.
          </h1>

          <div className="inline-flex items-center px-3.5 3xl:px-4 py-2 border-4 border-white/95 rounded-lg bg-[#d2e4f6]/90 text-[#17263d] text-4xl sm:text-5xl md:text-[56px] 3xl:text-[64px] font-matahari leading-none w-[90%]  md:w-[80%] lg:w-[90%] tracking-[-0.035em] shadow-[0_0_0_1px_rgba(126,177,229,0.65)] font-normal">
            Start owning it.
          </div>

          <p className="max-w-[490px] 3xl:max-w-[540px] text-white/92 text-sm sm:text-base 3xl:text-lg leading-[160%] tracking-[0.17px] font-light ">
            Convert tangible assets into digital tokens and open up new
            opportunities for global investment, fractional ownership, and
            borderless liquidity, across a trusted blockchain network.
          </p>

          <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 ">
            <a
              href="https://calendly.com/thomasenechi"
              className="px-6 py-3 rounded-lg font-medium text-center transition-all bg-primary hover:bg-primary-dark hover:-translate-y-0.5 text-white"
              target="_blank"
            >
              Tokenize an asset
            </a>
            <a
              href="/contact"
              className="px-6 py-3 rounded-lg font-medium text-center transition-all bg-white hover:bg-gray-100 hover:-translate-y-0.5 !text-black hover:!text-black"
            >
              Contact Us
            </a>
          </div>
        </div>

        {/* Hero Preview Cards */}
        <div className=" lg:flex items-end gap-1 3xl:gap-2 px-1 pt-12 md:pt-28 3xl:pt-32 -mr-5 md:-mr-24 xl:-mr-16 3xl:mr-0 max-w-full overflow-x-auto no-scrollbar md:overflow-visible  hidden ">
          {HERO_IMAGES.map((img, index) => (
            <button
              key={index}
              className={`relative flex-shrink-0 w-[170px] md:w-[200px] 3xl:w-[220px] h-[210px] md:h-[245px] 3xl:h-[270px] p-0 border-3 rounded-xl overflow-hidden cursor-pointer transition-all duration-300 transform origin-bottom hover:opacity-100 ${
                displayIndex === index
                  ? "border-white opacity-100 scale-y-[1.12]"
                  : "border-transparent opacity-80 scale-y-100"
              }`}
              onMouseEnter={() => setHoverIndex(index)}
              onMouseLeave={() => setHoverIndex(null)}
              onClick={() => handleCardClick(index)}
              type="button"
            >
              <Image
                src={img}
                alt={`Asset Preview ${index + 1}`}
                width={220}
                height={270}
                className="w-full h-full object-cover"
                priority={index === 0}
              />
            </button>
          ))}
        </div>
      </div>

      {/* Ecosystem Logos Bar — pushed to bottom via mt-auto in flex column */}
      <div className="relative z-5 mt-auto pb-1 pt-12 md:pt-2">
        <div className="page-container">
          <div className="flex items-center justify-center gap-2.5">
            <span className="h-[1px] flex-1 bg-white/18" />
            <p className="text-white/78 text-xs sm:text-sm font-semibold tracking-[0.045em] uppercase whitespace-nowrap ">
              REGULATED PARTICIPANTS IN OUR ECOSYSTEM
            </p>
            <span className="h-[1px] flex-1 bg-white/18" />
          </div>

          <div className="mt-2 flex items-center justify-between md:justify-center gap-6 md:gap-16 overflow-x-auto no-scrollbar pt-2 pb-4">
            <img
              src="/img/new-website/sec-Logo.png"
              alt="SEC Nigeria"
              className="max-h-[30px] md:max-h-[38px] w-auto object-contain flex-shrink-0"
            />
            <img
              src="/img/new-website/rmb.png"
              alt="RMB Nigeria Limited"
              className="max-h-[30px] md:max-h-[38px] w-auto object-contain flex-shrink-0"
            />
            <img
              src="/img/new-website/ARM-Logo.png"
              alt="Arm-logo"
              className="max-h-[30px] md:max-h-[38px] w-auto object-contain flex-shrink-0"
            />
            <img
              src="/img/new-website/BOI-TC-Logo.png"
              alt="BIO-TC Nigeria"
              className="max-h-[30px] md:max-h-[38px] w-auto object-contain flex-shrink-0"
            />
            <img
              src="/img/new-website/AO2law-logo.png"
              alt="AO2law-logo"
              className="max-h-[30px] md:max-h-[38px] w-auto object-contain flex-shrink-0"
            />
          </div>
        </div>
      </div>
    </section>
  );
};

export default Hero;

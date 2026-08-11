import React from "react";
import Image from "next/image";

const COMPLIANCE_LOGOS = [
  {
    icon: "/img/new-website/sec-logo.svg",
    desc: "SEC Nigeria RWA ",
  },

  {
    icon: "/img/new-website/rmb.svg",

    desc: "RMB Nigeria Limited",
  },
  {
    icon: "/img/new-website/sec-logo.svg",
    desc: "SEC Nigeria RWA ",
  },
  {
    icon: "/img/new-website/rmb.svg",

    desc: "RMB Nigeria Limited",
  },
];

const RegulatoryAlignment = () => {
  return (
    <section className="py-20 md:py-[120px] bg-white relative overflow-hidden">
      <div className=" mx-auto px-5 sm:px-8 md:px-12 lg:px-20 flex flex-col items-center text-center">
        {/* Header */}
        <div className="max-w-[750px] mx-auto mb-4 flex flex-col items-center">
          <p className="text-primary font-bold tracking-[0.15em] uppercase mb-4 text-xs sm:text-sm font-matahari">
            REGULATORY ALIGNMENT & COMPLIANCE
          </p>
          <h2 className="text-3xl sm:text-4xl md:text-5xl lg:text-[50px] font-extrabold leading-[1.15] mb-6 text-secondary font-matahari tracking-tight">
            Built With Regulatory Alignment
          </h2>
          <p className="text-gray-600 text-base sm:text-lg leading-[175%] font-sans font-normal max-w-[650px]">
            Trovotech is fully aligned with regulatory bodies and compliance
            standards, ensuring a secure and compliant ecosystem for tokenized
            real-world assets.
          </p>
        </div>

        {/* Logos Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 w-full  mx-auto  justify-center ">
          {COMPLIANCE_LOGOS.map((logo, idx) => (
            <div key={idx} className="flex  items-center text-center p-5 ">
              {/* Logo icon wrapper */}
              <div className="h-16 w-full flex items-center justify-center mb-4 relative">
                {logo.icon.endsWith(".svg") ? (
                  <img
                    src={logo.icon}
                    alt={logo.desc}
                    className="h-14 w-auto object-contain"
                    // className="h-14 w-auto object-contain filter grayscale hover:grayscale-0 hover:scale-105 transition-all duration-300"
                  />
                ) : (
                  <Image
                    src={logo.icon}
                    alt={logo.desc}
                    width={100}
                    height={56}
                    className="h-14 w-auto object-contain"
                    // className="h-14 w-auto object-contain filter grayscale hover:grayscale-0 hover:scale-105 transition-all duration-300"
                  />
                )}
              </div>
              <p className="text-gray-600 text-base sm:text-lg font-normal whitespace-nowrap">
                {logo.desc}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default RegulatoryAlignment;

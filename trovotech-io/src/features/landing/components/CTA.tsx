import React from "react";

const CTA = () => {
  return (
    <>
      <section className="py-20 md:py-[100px] bg-dark text-white text-center">
        <div className="page-container flex flex-col items-center">
          <h2 className="text-4xl sm:text-5xl md:text-[56px] 3xl:text-[64px] font-extrabold leading-[1.1] mb-8 font-matahari">
            Ready to access productive <br />
            capital infrastructure?
          </h2>

          <p className="text-lg sm:text-xl 3xl:text-[22px] text-gray-400 max-w-[600px] 3xl:max-w-[680px] leading-[160%] mb-12 font-sans">
            Whether you&apos;re an investor, asset sponsor, or institutional
            partner — Trovotech has a regulated pathway built for you.
          </p>

          <a
            href="https://calendly.com/thomasenechi"
            className="px-12 py-4 text-lg font-semibold bg-primary hover:bg-primary-dark text-white rounded-lg transition-all duration-300 hover:-translate-y-0.5 cursor-pointer"
            target="_blank"
          >
            Tokenize an Asset
          </a>
        </div>
      </section>

      {/* Decorative separating line */}
      <div className="bg-dark py-0">
        <div className="page-container">
          <div className="h-[1px] bg-white/28 w-full" />
        </div>
      </div>
    </>
  );
};

export default CTA;

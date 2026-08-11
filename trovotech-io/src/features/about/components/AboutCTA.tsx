import React from "react";

const AboutCTA = () => {
  return (
    <>
      <section className="py-20 md:py-32 bg-dark text-white relative overflow-hidden">
        <div className="relative z-10 max-w-[1280px] mx-auto px-5 sm:px-8 md:px-12 lg:px-20 flex flex-col items-center text-center">
          <h2 className="text-4xl sm:text-5xl md:text-[56px] font-extrabold leading-[1.1] mb-6 font-matahari max-w-[1000px]">
            Participate in the Future of Productive Asset Investment
          </h2>

          <p className="text-lg text-white/60 max-w-[600px] leading-[165%] mb-6 font-sans">
            Whether you are an investor seeking regulated access to real assets,
            an asset owner raising structured capital, or an institution ready
            to participate in new capital market infrastructure , there is a
            pathway built for you.
          </p>

          <div className="flex flex-col sm:flex-row items-center gap-4">
            <a
              href="/ecosystem"
              className="px-10 py-4 text-base font-semibold bg-primary hover:bg-primary-dark text-white rounded-lg transition-all duration-300 hover:-translate-y-0.5 cursor-pointer"
            >
              Join the Platform
            </a>
            <a
              href="/ecosystem"
              className="px-10 py-4 text-base font-semibold bg-transparent hover:bg-white/20 text-white border border-white rounded-lg transition-all duration-300 hover:-translate-y-0.5 cursor-pointer"
            >
              Explore Investments
            </a>
          </div>
          <p className="mt-5 text-white font-light">
            Speak with our team?{" "}
            <a
              href="/contact"
              className="text-white cursor-pointer font-semibold underline"
            >
              {" "}
              Contact us
            </a>{" "}
          </p>
        </div>
      </section>

      {/* Decorative separator for footer */}
      <div className="bg-dark py-0">
        <div className="max-w-[1280px] mx-auto px-5 sm:px-8 md:px-12 lg:px-20">
          <div className="h-[1px] bg-white/20 w-full" />
        </div>
      </div>
    </>
  );
};

export default AboutCTA;

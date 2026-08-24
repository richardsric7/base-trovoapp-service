import PrimaryButton from "@/component/PrimaryButton";

const AboutHero = () => {
  return (
    <div className="relative z-10 max-w-[1280px] mx-auto px-5 sm:px-8 md:px-12 lg:px-20 flex flex-col items-center text-center py-10 lg:py-8">
      <p className="text-primary font-bold tracking-[0.12em] uppercase mb-4 text-xs sm:text-base font-matahari">
        ABOUT TROVOTECH
      </p>
      <h1 className="text-3xl sm:text-4xl md:text-5xl lg:text-[54px] font-extrabold leading-[72px] mb-6 text-trovo-black max-w-[800px] font-matahari">
        Modernizing Access to Productive Asset Investment Across Africa.
      </h1>
      <p className="text-base  text-trovo-black max-w-[680px] leading-[170%] mb-4 font-sans">
        Trovotech is building regulated capital market infrastructure that
        enables transparent and structured investment into real-world assets
        through compliant digital financial systems.
      </p>

      <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mt-2 w-full sm:w-auto">
        <PrimaryButton
          label="Explore the platform"
          className="!text-white"
          href="/ecosystem"
        />
        <PrimaryButton
          label="Partner with us"
          className="!bg-transparent border !border-black !text-black hover:!bg-black/5"
          href="/contact"
        />
      </div>
    </div>
  );
};

export default AboutHero;

function TrovoBrand({ textColor }: { textColor: string }) {
  const textClasses = `text-primary-800 font-montserratMedium text-xl xl:text-2xl ${textColor}`;
  return (
    <div className="flex items-center space-x-3 mt-5 mb-10 px-3 w-full">
      <img className="w-10" src="/images/trovoLogo.png" alt="trovo logo" />
      <span className={textClasses}>Trovo App</span>
    </div>
  );
}

TrovoBrand.defaultProps = {
  textColor: 'text-primary-800',
};

export default TrovoBrand;

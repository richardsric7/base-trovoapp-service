import { useEffect, useState } from "react";

export const useCheckDeviceSize = () => {
  const [mobileDevice, setMobileDevice] = useState(false);
  useEffect(() => {
      const responsiveState = () => {
      setMobileDevice(window.innerWidth < 769);
    };
    window.addEventListener("resize", responsiveState);
    return ()=>{
        window.removeEventListener('resize', responsiveState)
    }
  }, [mobileDevice]);
  return { mobileDevice };
};

import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';

export default function Welcome2() {
  const navigate = useNavigate();
  return (
    <div className="flex flex-col h-3/4 lg:h-full space-y-3 lg:space-y-5 items-center justify-center ">
      <span className="rounded-full inline-block h-16" />
      <p className="text-primary-800 font-matahariExtended text-2xl md:text-4xl font-bold">
        Manage your Bantu
      </p>
      <p className="text-primary-700 font-matahariExtended text-2xl md:text-4xl font-bold">
        Digital Assets
      </p>
      <img
        className="px-5 md:px-0 w-500 h-500"
        src="/images/welcome2.png"
        alt="Welcome 2"
      />
      <div className="w-full h-12 flex items-center justify-center space-x-2">
        <span className="rounded-full inline-block w-3 h-3 bg-primary-800" />
        <span className="rounded-full inline-block w-3 h-3 bg-primary-800" />
        <span className="rounded-full inline-block w-3 h-3 bg-primary-500" />
      </div>
      <div className="w-3/4 md:w-1/4 flex space-x-2">
        <ButtonSecondary
          label="Previous"
          onclick={() => {
            navigate('/welcome1');
          }}
        />
        <Button
          label="Next"
          onclick={() => {
            navigate('/welcome3');
          }}
        />
      </div>
    </div>
  );
}

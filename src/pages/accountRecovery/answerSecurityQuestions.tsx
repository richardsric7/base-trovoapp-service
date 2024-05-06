import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TextInput from '../../components/textInput';
import Modal from '../../components/modal';
import TrovoBrand from '../../components/trovoBrand';
import { useFetchSecurityQuestionsQuery } from '../../store/api/authApi';
import { createAccount } from '../../utils/trovoSDK';
function AnswerSecurityQuestions() {
  const [showModal, setShowModal] = useState(false);
  const [tempAccount] = useState<Account>(createAccount());
  // const [fetchSecurityQuestions] = useFetchSecurityQuestionsQuery({
  //   signer: tempAccount.publicKey,
  //   publicKey: tempAccount.publicKey,
  //   secretKey: formInfo.secretKey,
  //   body: { userId: tempAccount.username, import: 1 },
  // });
  const navigate = useNavigate();

  return (
    <div className="flex h-screen items-center justify-center ">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Security
          </p>
          <p className="text-primary-700 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Questions
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="Welcome 1"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 h-full overflow-y-scroll">
        <div className="flex flex-col md:px-20 py-20 md:py-0 space-y-6 h-full items-center md:justify-center">
          <div className="flex flex-col space-y-6 md:h-full items-center justify-center">
            <p className="text-center w-full text-primary-800 text-2xl font-bold">
              Answer Security Questions
            </p>
            <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
              <p className="text-primary-800 text-md">
                Please answer the following security questions to proceed to the
                next step.
              </p>
            </div>
            <div className="w-3/4">
              <TextInput
                label="What is your mother’s maiden name?"
                leadingIcon="/images/question.png"
                inputType="text"
                placeholder="Enter answer"
                onInputChange={(newValue) => {
                  console.log('input has changed', newValue);
                }}
              />
            </div>
            <div className="w-3/4">
              <TextInput
                label="What is the name of the first street you lived at?"
                leadingIcon="/images/question.png"
                inputType="text"
                placeholder="Enter answer"
                onInputChange={(newValue) => {
                  console.log('input has changed', newValue);
                }}
              />
            </div>
            <div className="w-3/4">
              <TextInput
                label="What is your mother’s maiden name?"
                leadingIcon="/images/question.png"
                inputType="text"
                placeholder="Enter answer"
                onInputChange={(newValue) => {
                  console.log('input has changed', newValue);
                }}
              />
            </div>
            <Modal showModal={showModal} onClose={() => {}}>
              <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                <img src="/images/launch.png" alt="success" />
                <div className="flex flex-col text-center space-y-5 items-center w-2/3 mb-5 md:px-10 justify-center">
                  <p className="text-primary-800 text-md xl:text-lg font-semibold">
                    Congratulations!
                  </p>
                  <p className="text-primary-800 text-md xl:text-lg">
                    You have successfully recovered your account. Please copy
                    your secret key below to import your wallet afresh from your
                    device.
                  </p>
                </div>
                <div className="w-3/4 text-justify space-y-5 px-5 md:px-10 mb-5 py-5 bg-primary-100 rounded-xl">
                  <p className="text-primary-800 text-md xl:text-lg font-semibold">
                    Secret key for Ogbonge Wallet
                  </p>
                  <p className="text-primary-800 text-md break-all">
                    {/* eslint-disable-next-line */}
                    ASKONRWINOT-949I0IWRINEKLNFSKNFSDPG0J3-9JGRNKNGKN0J34INORGSDFW4W4WWEKLNDKLSFWKLNREKONELN4T4U48T53UONGKNGKLNK34T34WRKLGNKLGNREGKLNRKENGKLRNGEL
                  </p>
                  <div className="flex justify-between w-full">
                    <button type="button">
                      <img src="/images/copy.png" alt="copy" />
                    </button>
                    <button type="button">
                      <img src="/images/eyeShow.png" alt="show/hide" />
                    </button>
                  </div>
                </div>
                <div className="w-3/4">
                  <Button
                    label="Go to Dashboard"
                    onclick={() => {
                      navigate('/dashboard');
                      setShowModal(false);
                    }}
                  />
                </div>
              </div>
            </Modal>
            <div className="w-3/4">
              <Button
                label="Verify"
                onclick={() => {
                  // navigate('/register/backup');
                  setShowModal(true);
                }}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AnswerSecurityQuestions;

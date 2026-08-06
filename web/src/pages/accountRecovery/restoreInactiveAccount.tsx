import { useEffect, useState } from 'react';
// import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';
import TextInput from '../../components/textInput';
import Modal from '../../components/modal';
import TrovoBrand from '../../components/trovoBrand';
import { useSelector } from 'react-redux';
import {
  useFetchSecurityQuestionsQuery,
  useRestoreInactiveAccountMutation,
} from '../../store/api/authApi';
import { RootState } from '../../store/reduxStore';
import {
  showNotification,
  toggleLoader,
  showLoader,
  hideLoader,
} from '../../utils/showToaster';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import { useNavigate } from 'react-router-dom';
import Dropdown from '../../components/dropdown';

function RestoreInactiveAccount() {
  type SecurityQuestion = {
    id: number;
    question: string;
    answer: string;
    error: string;
  };

  type SecurityAnswers = {
    a1: string;
    a2: string;
    a3: string;
    q1: number;
    q2: number;
    q3: number;
  };

  const [showModal, setShowModal] = useState(false);
  const [showSecret, setShowSecret] = useState(false);
  const [errors, setErrors] = useState({
    a1: '',
    a2: '',
    a3: '',
    q1: '',
    q2: '',
    q3: '',
  });
  const [showEnsureBackupModal, setShowEnsureBackupModal] = useState(false);
  const [questions, setQuestions] = useState<SecurityQuestion[]>([]);
  const [securityAnswers, setSecurityAnswers] = useState<SecurityAnswers>({
    a1: '',
    a2: '',
    a3: '',
    q1: 0,
    q2: 0,
    q3: 0,
  });
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const [restoreInactiveAccount] = useRestoreInactiveAccountMutation();
  const navigate = useNavigate();

  const { data, isLoading } = useFetchSecurityQuestionsQuery({
    signer: tempData.publicKey,
    publicKey: tempData.publicKey,
    secretKey: tempData.secretKey,
    body: { username: tempData.username },
  });

  useEffect(() => {
    if (!tempData.emailOtp) {
      navigate('/recovery');
      return;
    }

    if (!isLoading) {
      hideLoader();
      const questions: SecurityQuestion[] = [];

      data.securityQuestions?.map((q: any) => {
        questions.push({
          id: q.ID,
          question: q.Question,
          answer: '',
          error: '',
        });
      });
      setQuestions(questions);
    } else {
      showLoader();
    }
  }, [isLoading]);

  const submitAnswers = async () => {
    const err = {
      a1: '',
      a2: '',
      a3: '',
      q1: '',
      q2: '',
      q3: '',
    };

    if (!securityAnswers.q1) {
      err.q1 = 'Please choose security question 1';
    }

    if (!securityAnswers.q2) {
      err.q2 = 'Please choose security question 2';
    }

    if (!securityAnswers.q3) {
      err.q3 = 'Please choose security question 3';
    }

    if (!securityAnswers.a1) {
      err.a1 = 'Please enter your answer';
    }

    if (!securityAnswers.a2) {
      err.a2 = 'Please enter your answer';
    }

    if (!securityAnswers.a3) {
      err.a3 = 'Please enter your answer';
    }

    setErrors({ ...err });
    // if there is no selected security question then there is no answer
    if (!securityAnswers.a1 || !securityAnswers.a2 || !securityAnswers.a3) {
      return;
    }

    try {
      toggleLoader();

      const body = {
        newSignerPublicKey: tempData.publicKey,
        emailOtp: tempData.emailOtp,
        username: tempData.username,
        securityAnswers,
      };

      const res = await restoreInactiveAccount({
        signer: tempData.publicKey,
        publicKey: tempData.publicKey,
        secretKey: tempData.secretKey,
        body,
      });

      toggleLoader();
      console.log('res', res);
      if ('data' in res) {
        setShowModal(true);
      } else if ('error' in res) {
        const errorResponse = res.error as ErrorResponse;
        showNotification(
          'error',
          errorResponse.data.message ??
            'Something went wrong. Please try again.',
          5000,
        );
      }
    } catch (error: any) {
      console.log(error);
      toggleLoader();
    }
  };

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
              Setup Security Questions
            </p>
            {questions.length > 0 ? (
              <form
                id="security-answers"
                onSubmit={(e: React.FormEvent<HTMLFormElement>) =>
                  e.preventDefault()
                }
                className="w-full space-y-6"
              >
                <div className="w-full text-center px-5 md:px-10 py-5 bg-primary-100 space-y-10 rounded-xl">
                  <p className="text-primary-800 text-md">
                    This will be required if you wish to make modifications to
                    your account in the future, and if you ever wish to opt in
                    for our account recovery service.
                  </p>
                  <p className="text-primary-800 text-md">
                    PLEASE DO NOT FORGET THE ANSWERS YOU PROVIDED FOR FUTURE
                    USE.
                  </p>
                </div>
                <div className="w-full">
                  <Dropdown
                    label="Choose question 1"
                    options={[
                      ...questions
                        .filter(
                          (q) =>
                            q.id !== securityAnswers.q2 &&
                            q.id !== securityAnswers.q3,
                        )
                        .map((q) => ({
                          text: q.question,
                          value: q.id,
                        })),
                    ]}
                    onSelect={(selectedItem) => {
                      setSecurityAnswers({
                        ...securityAnswers,
                        q1: selectedItem.value,
                      });
                      setErrors({ ...errors, q1: '', a1: '' });
                    }}
                  />
                  {errors.q1 && (
                    <p className="text-red-500 text-sm mt-1">{errors.q1}</p>
                  )}
                </div>
                {securityAnswers.q1 !== 0 && (
                  <div className="w-full">
                    <TextInput
                      label=""
                      leadingIcon="/images/question.png"
                      inputType="text"
                      placeholder="Enter answer"
                      defaultValue={securityAnswers.a1}
                      error={errors.a1}
                      onInputChange={(newValue) => {
                        securityAnswers.a1 = newValue;
                      }}
                    />
                  </div>
                )}
                <div className="w-full">
                  <Dropdown
                    label="Choose question 2"
                    options={[
                      ...questions
                        .filter(
                          (q) =>
                            q.id !== securityAnswers.q1 &&
                            q.id !== securityAnswers.q3,
                        )
                        .map((q) => ({
                          text: q.question,
                          value: q.id,
                        })),
                    ]}
                    onSelect={(selectedItem) => {
                      setSecurityAnswers({
                        ...securityAnswers,
                        q2: selectedItem.value,
                      });
                      setErrors({ ...errors, q2: '', a2: '' });
                    }}
                  />
                  {errors.q2 && (
                    <p className="text-red-500 text-sm mt-1">{errors.q2}</p>
                  )}
                </div>
                {securityAnswers.q2 !== 0 && (
                  <div className="w-full">
                    <TextInput
                      label=""
                      leadingIcon="/images/question.png"
                      inputType="text"
                      placeholder="Enter answer"
                      defaultValue={securityAnswers.a2}
                      error={errors.a2}
                      onInputChange={(newValue) => {
                        securityAnswers.a2 = newValue;
                      }}
                    />
                  </div>
                )}
                <div className="w-full">
                  <Dropdown
                    label="Choose question 3"
                    options={[
                      ...questions
                        .filter(
                          (q) =>
                            q.id !== securityAnswers.q1 &&
                            q.id !== securityAnswers.q2,
                        )
                        .map((q) => ({
                          text: q.question,
                          value: q.id,
                        })),
                    ]}
                    onSelect={(selectedItem) => {
                      setSecurityAnswers({
                        ...securityAnswers,
                        q3: selectedItem.value,
                      });
                      setErrors({ ...errors, q3: '', a3: '' });
                    }}
                  />
                  {errors.q3 && (
                    <p className="text-red-500 text-sm mt-1">{errors.q3}</p>
                  )}
                </div>
                {securityAnswers.q3 !== 0 && (
                  <div className="w-full">
                    <TextInput
                      label=""
                      leadingIcon="/images/question.png"
                      inputType="text"
                      placeholder="Enter answer"
                      defaultValue={securityAnswers.a3}
                      error={errors.a3}
                      onInputChange={(newValue) => {
                        securityAnswers.a3 = newValue;
                      }}
                    />
                  </div>
                )}
                <div className="w-full">
                  <Button
                    type="button"
                    label="Submit"
                    onclick={submitAnswers}
                  />
                </div>
              </form>
            ) : (
              <div className="space-y-6">
                <div className="w-full text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
                  <p className="text-primary-800 text-md">
                    You currently have no security questions setup for this
                    account yet. Please setup your security questions and then
                    enable account recovery.
                  </p>
                </div>
                <Button
                  label="Go Back"
                  onclick={() => {
                    navigate(-1);
                  }}
                />
              </div>
            )}
            <Modal
              showModal={showModal}
              onClose={() => {
                setShowModal(false);
              }}
            >
              <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                <img src="/images/launch.png" alt="success" />
                <div className="flex flex-col text-center space-y-5 items-center w-2/3 mb-5 md:px-10 justify-center">
                  <p className="text-primary-800 text-md xl:text-lg font-semibold">
                    Congratulations!
                  </p>
                  <p className="text-primary-800 text-md xl:text-lg">
                    We have generated new credentials for your account. Please
                    copy your new secret key below to import your wallet afresh
                    from your device.
                  </p>
                </div>
                <div className="w-3/4 text-justify space-y-5 px-5 md:px-10 mb-5 py-5 bg-primary-100 rounded-xl">
                  <div className="flex flex-col space-y-5">
                    <p className="text-left w-full text-primary-800 text-md font-bold">
                      Alias:
                    </p>
                    <div className="flex justify-between">
                      <p>{tempData.username}</p>
                      <button
                        type="button"
                        onClick={() =>
                          navigator.clipboard
                            .writeText(tempData.username)
                            .then(() => {
                              showNotification('info', 'Username copied!');
                            })
                        }
                      >
                        <img src="/images/copy.png" alt="copy" />
                      </button>
                    </div>
                  </div>
                  <div className="flex flex-col space-y-5 w-full">
                    <p className="text-left w-full text-primary-800 text-md font-bold">
                      Public Key:
                    </p>
                    <div className="flex w-full justify-between">
                      <div className="w-4/5 h-full break-all">
                        {tempData.publicKey}
                      </div>
                      <button
                        type="button"
                        onClick={() =>
                          navigator.clipboard
                            .writeText(tempData.publicKey)
                            .then(() => {
                              showNotification('info', 'Public key copied!');
                            })
                        }
                      >
                        <img src="/images/copy.png" alt="copy" />
                      </button>
                    </div>
                  </div>
                  <p className="text-primary-800 text-md xl:text-lg font-semibold">
                    Secret key
                  </p>
                  <div className="flex justify-between w-full">
                    <p className="w-4/5 text-md break-all">
                      {/* eslint-disable-next-line */}
                      {showSecret ? tempData.secretKey : '***********'}
                    </p>
                    <div className="flex justify-between space-x-5">
                      <button
                        onClick={async () => {
                          navigator.clipboard
                            .writeText(tempData.secretKey)
                            .then(() => {
                              showNotification('info', 'Secret key copied!');
                            });
                        }}
                        type="button"
                      >
                        <img src="/images/copy.png" alt="copy" />
                      </button>
                      <button
                        onClick={async () => {
                          setShowSecret(!showSecret);
                        }}
                        type="button"
                      >
                        <img src="/images/eyeShow.png" alt="show/hide" />
                      </button>
                    </div>
                  </div>
                </div>
                <div className="w-3/4">
                  <Button
                    label="Continue"
                    onclick={() => {
                      setShowModal(false);
                      setShowEnsureBackupModal(true);
                    }}
                  />
                </div>
              </div>
            </Modal>
            <Modal
              showModal={showEnsureBackupModal}
              onClose={() => {
                setShowEnsureBackupModal(false);
              }}
            >
              <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                <img src="/images/launch.png" alt="success" />
                <div className="flex flex-col text-center space-y-5 items-center w-3/4 mb-5 md:px-10 justify-center">
                  <p className="text-primary-800 text-md xl:text-xl font-bold">
                    Account Recovery
                  </p>
                  <>
                    <p className="text-primary-800 text-md xl:text-lg">
                      Your account has been successfully recovered. Before
                      proceeding to import your account please confirm that you
                      have backed up the new account information. If you have
                      not backed it up, kindly tap the back button and back it
                      up.
                    </p>
                    <ButtonSecondary
                      label="Go back and backup"
                      onclick={() => {
                        setShowModal(true);
                        setShowEnsureBackupModal(false);
                      }}
                    />
                    <Button
                      label="Continue, I have backed up"
                      onclick={() => {
                        setShowEnsureBackupModal(false);
                        showNotification(
                          'success',
                          'Your account has successfully been recovered. You can now import your account with the new secret key.',
                          5000, // delay for 5secs
                        );
                        navigate('/import');
                      }}
                    />
                  </>
                </div>
              </div>
            </Modal>
          </div>
        </div>
      </div>
    </div>
  );
}

export default RestoreInactiveAccount;

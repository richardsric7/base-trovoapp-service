import { useEffect, useState } from 'react';
// import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';
import TextInput from '../../components/textInput';
import Modal from '../../components/modal';
import TrovoBrand from '../../components/trovoBrand';
import { useDispatch, useSelector } from 'react-redux';
import { setFormState } from '../../store/authSlice';
import {
  useFetchSecurityQuestionsQuery,
  useSubmitSecurityAnswersMutation,
  useRequestAccountRecoveryMutation,
  useLazyGetRecoveryStatusQuery,
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

function AnswerSecurityQuestions() {
  type SecurityQuestion = {
    id: number;
    question: string;
    answer: string;
    error: string;
  };

  const [showModal, setShowModal] = useState(false);
  const [showSecret, setShowSecret] = useState(false);
  const [messages, setMessages] = useState([]);
  const [backupDone, setBackupDone] = useState(false);
  const [showEnsureBackupModal, setShowEnsureBackupModal] = useState(false);
  // when the started recovery takes effect (set once it has started)
  const [pendingUntil, setPendingUntil] = useState<Date | null>(null);
  const [questions, setQuestions] = useState<SecurityQuestion[]>([]);
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const [requestAccountRecovery] = useRequestAccountRecoveryMutation();
  const [submitSecurityAnswers] = useSubmitSecurityAnswersMutation();
  const [getRecoveryStatus] = useLazyGetRecoveryStatusQuery();
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const clearRecoveryData = () => {
    dispatch(
      setFormState({
        ...tempData,
        emailOtp: '',
        username: '',
        secretKey: '',
        address: '',
      }),
    );
  };

  // a started recovery completes after the recovery period, unless the
  // account owner cancels it with their current key: follow it
  useEffect(() => {
    if (!pendingUntil) return undefined;
    const check = async () => {
      const res = await getRecoveryStatus({
        signer: tempData.address,
        address: tempData.address,
        secretKey: tempData.secretKey,
        body: { username: tempData.username },
      });
      const status = res.data?.status;
      if (status === 'COMPLETED') {
        clearRecoveryData();
        navigate('/import');
        showNotification(
          'success',
          'Your account has been recovered. You can now import it with the new secret key.',
          5000,
        );
      } else if (status === 'CANCELED' || status === 'FAILED') {
        clearRecoveryData();
        navigate('/login');
        showNotification(
          'error',
          status === 'CANCELED'
            ? 'The recovery was canceled from a device that has the current key.'
            : 'The recovery could not be completed. Please contact support.',
          8000,
        );
      }
    };
    check();
    const timer = setInterval(check, 30000);
    return () => clearInterval(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pendingUntil]);

  const { data, isLoading } = useFetchSecurityQuestionsQuery({
    signer: tempData.address,
    address: tempData.address,
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
      const securityAnswers: {
        a1: '';
        a2: '';
        a3: '';
        id: number;
        q1: number;
        q2: number;
        q3: number;
      } = data?.userSecurityAnswers;
      const questions: SecurityQuestion[] = [];

      data.securityQuestions?.map((q: any) => {
        if (
          q.ID === securityAnswers.q1 ||
          q.ID === securityAnswers.q2 ||
          q.ID === securityAnswers.q3
        ) {
          questions.push({
            id: q.ID,
            question: q.Question,
            answer: '',
            error: '',
          });
        }
      });
      setQuestions(questions);
    } else {
      showLoader();
    }
  }, [isLoading]);

  const submitAnswers = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    let hasError = false;
    for (const q of questions) {
      if (!q.answer) {
        showNotification('error', 'Please enter all your security answers!');
        q.error = 'Please enter your answer';
        hasError = true;
      } else {
        q.error = '';
      }
    }

    if (hasError) return;

    try {
      toggleLoader();

      const data = {
        q1: questions[0].id,
        a1: questions[0].answer,
        q2: questions[1].id,
        a2: questions[1].answer,
        q3: questions[2].id,
        a3: questions[2].answer,
      };

      const res = await submitSecurityAnswers({
        signer: tempData.address,
        address: tempData.address,
        secretKey: tempData.secretKey,
        body: { username: tempData.username, answers: data },
      });

      toggleLoader();
      console.log('res', res);
      if ('data' in res) {
        // navigate('/backup');
        setShowModal(true);
      } else if ('error' in res) {
        const errorResponse = res.error as ErrorResponse;
        showNotification(
          'error',
          errorResponse.data.error ?? 'Something went wrong. Please try again.',
        );
      }
    } catch (error: any) {
      console.log(error);
      toggleLoader();
    }
  };

  const submitRequestAccountRecovery = async (commit: number) => {
    try {
      toggleLoader();

      const body = {
        newSignerAddress: tempData.address,
        commit,
        emailOtp: tempData.emailOtp,
        username: tempData.username,
        transactionId: '',
        securityAnswers: {
          q1: questions[0].id,
          a1: questions[0].answer,
          q2: questions[1].id,
          a2: questions[1].answer,
          q3: questions[2].id,
          a3: questions[2].answer,
        },
      };

      const res = await requestAccountRecovery({
        signer: tempData.address,
        address: tempData.address,
        secretKey: tempData.secretKey,
        body,
      });

      toggleLoader();
      console.log('res', res);
      if ('data' in res) {
        setShowModal(false);
        setMessages(res.data.messages);
        if (commit === 0) {
          setShowEnsureBackupModal(true);
        } else {
          setShowEnsureBackupModal(false);
          setPendingUntil(
            res.data.executeAfter ? new Date(res.data.executeAfter) : new Date(),
          );
        }
      } else if ('error' in res) {
        const errorResponse = res.error as ErrorResponse;
        showNotification(
          'error',
          errorResponse.data.error ?? 'Something went wrong. Please try again.',
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
              Answer Security Questions
            </p>
            {pendingUntil ? (
              <div className="w-3/4 space-y-6">
                <div className="w-full text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl space-y-3">
                  <p className="text-primary-800 text-md font-semibold">
                    Recovery in progress
                  </p>
                  <p className="text-primary-800 text-md">
                    Your wallets move to your new key on{' '}
                    {pendingUntil.toLocaleString()}. Until then the account
                    owner is notified and can cancel the recovery from a device
                    that still has the current key.
                  </p>
                  <p className="text-primary-800 text-md">
                    Keep your new secret key safe. You can leave this page and
                    import your account with the new secret key after that
                    time; this page checks the status every 30 seconds.
                  </p>
                </div>
                {messages.map((m, index) => (
                  <p
                    className="text-primary-800 text-sm text-center"
                    key={`pending-${index}`}
                  >
                    {m}
                  </p>
                ))}
              </div>
            ) : questions.length > 0 ? (
              <form
                id="security-answers"
                onSubmit={submitAnswers}
                className="w-3/4 space-y-6"
              >
                <div className="w-full text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
                  <p className="text-primary-800 text-md">
                    Please answer the following security questions to proceed to
                    the next step.
                  </p>
                </div>
                {questions.map((q) => (
                  <div
                    className="w-full"
                    key={`${new Date().getTime()}${q.id}`}
                  >
                    <TextInput
                      label={q.question}
                      leadingIcon="/images/question.png"
                      inputType="text"
                      placeholder="Enter answer"
                      defaultValue={q.answer}
                      error={q.error}
                      onInputChange={(newValue) => {
                        q.answer = newValue;
                      }}
                    />
                  </div>
                ))}
                <div className="w-full">
                  <Button
                    type="submit"
                    label="Verify"
                    onclick={() => {
                      console.log(questions);
                    }}
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
                      Address:
                    </p>
                    <div className="flex w-full justify-between">
                      <div className="w-4/5 h-full break-all">
                        {tempData.address}
                      </div>
                      <button
                        type="button"
                        onClick={() =>
                          navigator.clipboard
                            .writeText(tempData.address)
                            .then(() => {
                              showNotification('info', 'Address copied!');
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
                      submitRequestAccountRecovery(0); // 0 = dry run
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
                  {!backupDone && (
                    <>
                      <p className="text-primary-800 text-md xl:text-lg">
                        Before completing account recovery please confirm that
                        you have backed up the new account information. If you
                        have not backed it up, kindly tap the back button and
                        back it up.
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
                          // setShowModal(false);
                          setBackupDone(!backupDone);
                        }}
                      />
                    </>
                  )}
                  {backupDone && (
                    <>
                      {messages.map((m, index) => (
                        <p
                          className="text-danger text-md xl:text-lg"
                          key={`${new Date().getTime()}${index}`}
                        >
                          {m}
                        </p>
                      ))}
                      <Button
                        label="Ok, start account recovery"
                        onclick={() => {
                          setShowModal(false);
                          submitRequestAccountRecovery(1); // 1 = final commit
                        }}
                      />
                      <ButtonSecondary
                        label="No, cancel"
                        onclick={() => {
                          setShowModal(false);
                          setBackupDone(false);
                          setShowEnsureBackupModal(false);
                        }}
                      />
                    </>
                  )}
                </div>
              </div>
            </Modal>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AnswerSecurityQuestions;

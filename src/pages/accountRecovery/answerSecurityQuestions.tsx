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
} from '../../store/api/authApi';
import { RootState } from '../../store/reduxStore';
import { showNotification, toggleLoader } from '../../utils/showToaster';
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
  const [invalidateOldSigner, setInvalidateOldSigner] = useState(false);
  const [questions, setQuestions] = useState<SecurityQuestion[]>([]);
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const [requestAccountRecovery] = useRequestAccountRecoveryMutation();
  const [submitSecurityAnswers] = useSubmitSecurityAnswersMutation();
  const navigate = useNavigate();
  const dispatch = useDispatch();

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
        signer: tempData.publicKey,
        publicKey: tempData.publicKey,
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
        newSignerPublicKey: tempData.publicKey,
        disableOldSignerFromPrimaryWallet: invalidateOldSigner ? 1 : 0,
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
        signer: tempData.publicKey,
        publicKey: tempData.publicKey,
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
          navigate('/import');
          showNotification(
            'success',
            'Your account has successfully been recovered. You can now import your account with the new secret key.',
            5000, // delay for 5secs
          );
          dispatch(
            setFormState({
              ...tempData,
              emailOtp: '',
              username: '',
              secretKey: '',
              publicKey: '',
            }),
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
            <div className="w-3/4 text-center px-5 md:px-10 py-5 bg-primary-100 rounded-xl">
              <p className="text-primary-800 text-md">
                Please answer the following security questions to proceed to the
                next step.
              </p>
            </div>
            <form
              id="security-answers"
              onSubmit={submitAnswers}
              className="w-3/4 space-y-6"
            >
              {questions &&
                questions.map((q) => (
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
                <div className="w-3/4 flex justify-center space-x-3">
                  <input
                    type="checkbox"
                    name="import"
                    defaultChecked={invalidateOldSigner}
                    onChange={() => {
                      setInvalidateOldSigner(!invalidateOldSigner);
                    }}
                  />
                  <p className="text-gray-500">
                    Invalidate old signer from primary wallet?
                  </p>
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
                        label="Ok, complete account recovery"
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

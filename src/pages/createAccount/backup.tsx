import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/button';
import TrovoBrand from '../../components/trovoBrand';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import { Encryptor } from '../../utils/encryptor';
import { showNotification } from '../../utils/showToaster';
import { PasswordInputModal } from '../../components/passwordModal';

export default function Backup() {
  const [showCreds, setShowCredentials] = useState(false);
  const [showEnterPassword, setShowEnterPassword] = useState(false);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const [password, setPassword] = useState(tempData.password ?? '');
  const [secret, setSecret] = useState('');
  const navigate = useNavigate();

  function ShowCredentials() {
    return (
      <div className="text-primary-800 space-y-5 text-md w-full text-justify">
        <p>
          Please copy the following details correctly and store in a safe place
        </p>
        <p>You can copy your secret key as displayed below:</p>
        <div className="flex flex-col space-y-5">
          <p className="text-left w-full text-primary-800 text-md font-bold">
            Alias:
          </p>
          <div className="flex justify-between">
            <p>{appUser.username}</p>
            <button
              type="button"
              onClick={() =>
                navigator.clipboard.writeText(appUser.username).then(() => {
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
            <div className="w-4/5 h-full break-all">{appUser.publicKey}</div>
            <button
              type="button"
              onClick={() =>
                navigator.clipboard.writeText(appUser.publicKey).then(() => {
                  showNotification('info', 'Public key copied!');
                })
              }
            >
              <img src="/images/copy.png" alt="copy" />
            </button>
          </div>
        </div>
        <div className="flex flex-col space-y-5 w-full">
          <p className="text-left w-full text-primary-800 text-md font-bold">
            Secret Key:
          </p>
          <div className="flex w-full justify-between">
            <div className="w-4/5 h-full break-all">
              {secret ? secret : '***********'}
            </div>
            <div className="flex justify-between space-x-5">
              <button
                type="button"
                onClick={async () => {
                  if (!secret) {
                    const encryptor = new Encryptor();
                    const decryptedData = await encryptor.decryptData(
                      appUser.secretKeys[0],
                      password,
                      appUser.primarySigner,
                    );

                    setSecret(decryptedData);
                  } else {
                    setSecret('');
                  }
                }}
              >
                <img src="/images/eyeShow.png" alt="copy" />
              </button>
              <button
                type="button"
                onClick={async () => {
                  const encryptor = new Encryptor();
                  const decryptedData = await encryptor.decryptData(
                    appUser.secretKeys[0],
                    password,
                    appUser.primarySigner,
                  );
                  navigator.clipboard.writeText(decryptedData).then(() => {
                    showNotification('info', 'Secret key copied!');
                  });
                }}
              >
                <img src="/images/copy.png" alt="copy" />
              </button>
            </div>
          </div>
        </div>
        <div className="flex justify-end w-full">
          <button
            type="button"
            onClick={async () => {
              const encryptor = new Encryptor();
              const decryptedData = await encryptor.decryptData(
                appUser.secretKeys[0],
                password,
                appUser.primarySigner,
              );

              navigator.clipboard
                .writeText(
                  `Wallet alias: ${appUser.username}\nPublic key: ${appUser.publicKey}\nSecret key: ${decryptedData}`,
                )
                .then(() => {
                  showNotification('info', 'Wallet info copied!');
                });
            }}
          >
            Copy All
          </button>
        </div>
        <div className="w-3/4">
          <Button
            label="Continue"
            onclick={() => {
              navigate('/dashboard');
            }}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="flex md:h-full items-center justify-center">
      <div className="hidden md:block w-3/5 h-full p-3">
        <div className="flex h-full space-y-3 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <TrovoBrand />
          <p className="text-primary-800 font-matahariExtended text-center text-3xl xl:text-4xl font-bold">
            Backup
          </p>
          <img
            className="w-100 h-100"
            src="/images/OTPsecurity.png"
            alt="verification"
          />
        </div>
      </div>
      <div className="w-full md:w-2/5 pt-20 lg:pt-0">
        <div className="flex flex-col px-5 md:px-20 space-y-6 md:h-full items-center justify-center">
          <p className="text-left w-full text-primary-800 text-2xl font-bold">
            Backup Wallet
          </p>
          {showCreds ? (
            <ShowCredentials />
          ) : (
            <ShowPreliminary
              value={showCreds}
              onInputChange={(value) => {
                console.log('changed to: ', value);
                if (password) {
                  setShowCredentials(true);
                  return;
                }

                setShowCredentials(false);
                setShowEnterPassword(true);
              }}
            />
          )}
          <PasswordInputModal
            show={showEnterPassword}
            onClose={() => {
              setShowEnterPassword(false);
            }}
            onDone={(value) => {
              setPassword(value);
              setShowEnterPassword(false);
              setShowCredentials(true);
            }}
          />
        </div>
      </div>
    </div>
  );
}

function ShowPreliminary({
  onInputChange,
  value,
}: {
  onInputChange: (value: boolean) => void;
  value: boolean;
}) {
  return (
    <div className="text-primary-800 space-y-5 text-md text-justify">
      <p>
        I have ensured that no one is looking and I understand that I should
        never share my secret key with anyone.
      </p>
      <p>
        I understand that I need to securely store my secret key and that if
        this app is deleted or moved to another device without opting in for the
        account recovery service, I can only restore my wallet with the secret
        key.
      </p>
      <p>
        I understand that if I lose my secret key, Trovotech is not liable to
        any loss and my funds are securely held and controlled on this device
        not by Trovotech.
      </p>
      <div className="w-full flex space-x-3">
        <input
          type="checkbox"
          name="import"
          checked={value}
          onChange={(e) => {
            onInputChange(e.target.value as unknown as boolean);
          }}
        />
        <p className="text-gray-500">
          I have read and understood all the above
        </p>
      </div>
    </div>
  );
}

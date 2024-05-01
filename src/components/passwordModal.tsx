import { useState } from 'react';
import Modal from './modal';
import TextInput from './textInput';
import Button from './button';
import { Encryptor } from '../utils/encryptor';
import { useSelector } from 'react-redux';
import { RootState } from '../store/reduxStore';

type Props = {
  show: boolean;
  onClose: () => void;
  onDone: (password: string) => void;
};

export function PasswordInputModal({ show, onClose, onDone }: Props) {
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const [password, setPassword] = useState('');
  const [passwordErr, setPasswordErr] = useState('');

  const isValidPassword = async () => {
    try {
      const encryptor = new Encryptor();
      await encryptor.decryptData(
        appUser.secretKeys[0],
        password,
        appUser.publicKey,
      );
      return true;
    } catch (error: any) {
      return false;
    }
  };
  return (
    <Modal
      showModal={show}
      onClose={() => {
        onClose();
      }}
    >
      <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
        <p className="text-center w-full text-primary-800 text-2xl font-bold">
          Authorize
        </p>
        <img
          className="w-36"
          src="/images/welcome.png"
          alt="Enter password to authorize"
        />
        <div className="w-3/4 space-y-1">
          <TextInput
            label=""
            leadingIcon="/images/lock.png"
            inputType="password"
            placeholder="Enter password"
            onInputChange={(newValue) => {
              setPassword(newValue);
            }}
          />
          {passwordErr && <p className="text-red-500 text-sm">{passwordErr}</p>}
        </div>
        <div className="w-3/4">
          <Button
            label="Ok"
            onclick={async () => {
              setPasswordErr('');

              if (!password) {
                setPasswordErr('Please enter a password!');
                return;
              }

              if (!(await isValidPassword())) {
                setPasswordErr('Password is invalid!');
                return;
              }

              onDone(password);
            }}
          />
        </div>
      </div>
    </Modal>
  );
}

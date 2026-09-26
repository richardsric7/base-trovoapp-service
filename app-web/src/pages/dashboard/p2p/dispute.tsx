import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Header from '../../../components/header';
import { useOpenP2PDisputeMutation } from '../../../store/api/p2pApis';
import { useP2PIdentity } from '../../../hooks/useP2PIdentity';
import { p2p } from '../../../components/p2p/P2PTheme';

const SUBJECTS: Record<string, string> = {
  NO_PAYMENT: 'No payment was made',
  UNDER_PAYMENT: 'Underpayment',
  OVER_PAYMENT: 'Overpayment',
  PAYMENT_NOT_RECEIVED: 'Payment not received',
  PAYMENT_MARKED_SENT_IN_ERROR: 'Payment marked sent in error',
  WRONG_PAYMENT_AMOUNT: 'Wrong payment amount',
  OTHER: 'Other',
};

// P2PDispute (Plan Section 96.9): file a dispute (subject + description).
// No live chat - the legacy schema's order_chat_messages table was
// explicitly decided to be irrelevant prior art, not something to add here.
export default function P2PDispute() {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const { creds, ready } = useP2PIdentity();
  const [openDispute] = useOpenP2PDisputeMutation();
  const [subject, setSubject] = useState('PAYMENT_NOT_RECEIVED');
  const [description, setDescription] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const submit = async () => {
    if (!ready || !orderId) return;
    if (description.trim().length < 10) {
      setError('Please add a few more details.');
      return;
    }
    setSubmitting(true);
    setError('');
    try {
      await openDispute({ creds, orderId, subject, description }).unwrap();
      navigate(`/dashboard/p2p/order/${orderId}`);
    } catch (e: any) {
      setError(e?.data?.message || 'Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col space-y-4 p-3 max-w-2xl">
      <Header />
      <h1 className="text-xl font-bold text-primary-800 px-2">Raise a dispute</h1>

      <div>
        <label className="font-semibold block mb-2">What went wrong?</label>
        <select
          value={subject}
          onChange={(e) => setSubject(e.target.value)}
          className="w-full p-3 rounded-xl bg-white outline-none"
        >
          {Object.entries(SUBJECTS).map(([k, v]) => (
            <option key={k} value={k}>
              {v}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="font-semibold block mb-2">Describe what happened</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={5}
          placeholder="What happened, and what would resolve it?"
          className="w-full p-3 rounded-xl bg-white outline-none"
        />
      </div>

      {error && <div className="text-red-600">{error}</div>}

      <button
        disabled={submitting}
        onClick={submit}
        className="w-full py-4 rounded-2xl font-semibold text-white disabled:opacity-50"
        style={{ backgroundColor: p2p.danger }}
      >
        Submit dispute
      </button>
    </div>
  );
}

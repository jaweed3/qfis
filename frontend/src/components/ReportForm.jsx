import { useState, useCallback } from 'react';
import { submitReport } from '../lib/api';

export default function ReportForm({ onSubmitted }) {
  const [merchantId, setMerchantId] = useState('');
  const [note, setNote] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const [reportId, setReportId] = useState('');
  const [error, setError] = useState(null);

  const handleSubmit = useCallback(async () => {
    if (!merchantId.trim()) return;
    try {
      const res = await submitReport(merchantId.trim(), note.trim());
      setReportId(res.report_id);
      setSubmitted(true);
      setError(null);
      if (onSubmitted) onSubmitted(res);
      setTimeout(() => {
        setSubmitted(false);
        setMerchantId('');
        setNote('');
        setReportId('');
      }, 3000);
    } catch (e) {
      setError(e.message);
    }
  }, [merchantId, note, onSubmitted]);

  return (
    <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg p-5">
      <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3 tracking-widest">
        LAPORKAN TEMUAN
      </div>

      {submitted ? (
        <div className="p-4 rounded-md text-center font-mono text-[13px]" style={{ background: 'rgba(48,209,88,0.12)', border: '1px solid rgba(48,209,88,0.25)', color: '#30D158' }}>
          ✓ Laporan diterima. ID: {reportId}
        </div>
      ) : (
        <div className="flex flex-col gap-2.5">
          <input
            value={merchantId}
            onChange={(e) => setMerchantId(e.target.value)}
            placeholder="Merchant ID atau nama..."
            className="w-full px-4 py-[11px] rounded-md font-mono text-[13px] outline-none"
            style={{
              background: '#07070F',
              border: '1px solid rgba(43,43,255,0.2)',
              color: '#FFFFFF',
            }}
          />
          <textarea
            value={note}
            onChange={(e) => setNote(e.target.value)}
            placeholder="Catatan tambahan (opsional)..."
            rows={3}
            className="w-full px-4 py-[11px] rounded-md font-body text-[13px] outline-none resize-none"
            style={{
              background: '#07070F',
              border: '1px solid rgba(43,43,255,0.2)',
              color: '#FFFFFF',
            }}
          />
          {error && <div className="text-xs font-mono text-[#FF2D55]">{error}</div>}
          <button
            onClick={handleSubmit}
            className="w-full font-mono text-[11px] font-bold uppercase tracking-wider py-2.5 rounded-md cursor-pointer transition-all"
            style={{
              background: 'rgba(255,159,10,0.12)',
              color: '#FF9F0A',
              border: '1px solid rgba(255,159,10,0.3)',
            }}
          >
            KIRIM LAPORAN →
          </button>
        </div>
      )}
    </div>
  );
}

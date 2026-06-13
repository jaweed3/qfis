import { useState, useCallback } from 'react';
import RiskBadge, { riskColor } from './RiskBadge';
import { checkMerchant } from '../lib/api';

export default function MerchantChecker() {
  const [input, setInput] = useState('');
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleCheck = useCallback(async () => {
    if (!input.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const data = await checkMerchant(input.trim());
      setResult(data);
    } catch (e) {
      setError(e.message);
      setResult(null);
    }
    setLoading(false);
  }, [input]);

  return (
    <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg p-5">
      <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3 tracking-widest">
        CEK MERCHANT
      </div>
      <div className="flex gap-2">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleCheck()}
          placeholder="QRIS-XXXX atau nama merchant..."
          className="flex-1 px-4 py-[11px] rounded-md font-mono text-[13px] outline-none"
          style={{
            background: '#07070F',
            border: '1px solid rgba(43,43,255,0.2)',
            color: '#FFFFFF',
          }}
        />
        <button
          onClick={handleCheck}
          disabled={loading}
          className="font-mono text-[11px] font-bold uppercase tracking-wider px-6 py-[11px] rounded-md border-none cursor-pointer transition-all disabled:opacity-50"
          style={{
            background: loading ? 'rgba(43,43,255,0.2)' : '#2B2BFF',
            color: loading ? 'rgba(255,255,255,0.4)' : '#FFFFFF',
          }}
        >
          {loading ? 'SCAN...' : 'SCAN →'}
        </button>
      </div>

      {error && (
        <div className="mt-4 p-3 rounded-md text-xs font-mono" style={{ background: 'rgba(255,45,85,0.12)', border: '1px solid rgba(255,45,85,0.3)', color: '#FF2D55' }}>
          {error}
        </div>
      )}

      {result && !loading && !error && (
        <div
          className="mt-4 p-4 rounded-md"
          style={{ background: '#07070F', border: `1px solid ${riskColor(result.risk_score)}40` }}
        >
          <div className="flex justify-between items-start">
            <div>
              <div className="font-mono text-[13px] text-white font-semibold">{input}</div>
              <div className="font-mono text-[11px] text-[rgba(255,255,255,0.4)] mt-0.5">{result.flags.length} flag(s)</div>
            </div>
            <RiskBadge score={result.risk_score} />
          </div>

          <div className="mt-4 text-xs font-mono text-[rgba(255,255,255,0.55)] leading-relaxed">
            {result.recommendation}
          </div>

          {result.flags.length > 0 && (
            <div className="mt-3 space-y-1">
              {result.flags.map((f, i) => (
                <div key={i} className="text-[11px] font-mono" style={{ color: '#FF9F0A' }}>
                  ⚠ {f}
                </div>
              ))}
            </div>
          )}

          {result.risk_score >= 70 && (
            <div className="mt-3 p-2 rounded-md text-xs font-mono" style={{ background: 'rgba(255,45,85,0.12)', border: '1px solid rgba(255,45,85,0.3)', color: '#FF2D55' }}>
              ⚠ Merchant ini terindikasi aktivitas mencurigakan. Siap dilaporkan ke OJK.
            </div>
          )}
        </div>
      )}
    </div>
  );
}

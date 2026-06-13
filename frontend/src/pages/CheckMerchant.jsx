import { useState, useEffect } from 'react';
import MerchantChecker from '../components/MerchantChecker';
import { getMerchantNetwork } from '../lib/api';
import { riskColor } from '../components/RiskBadge';

export default function CheckMerchant() {
  const [merchants, setMerchants] = useState([]);

  useEffect(() => {
    getMerchantNetwork().then(setMerchants).catch(() => {});
  }, []);

  return (
    <div className="max-w-[600px] mx-auto animate-[fade-up_0.4s_ease-out]">
      <div className="mb-6">
        <div className="font-display text-2xl font-bold text-white mb-2">Cek Merchant QRIS</div>
        <div className="font-body text-[15px] text-[rgba(255,255,255,0.55)]">
          Masukkan QRIS merchant ID atau nama merchant untuk mendapatkan risk assessment instan.
        </div>
      </div>

      <MerchantChecker />

      {merchants.length > 0 && (
        <div className="mt-4 p-5 bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg">
          <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3 tracking-widest">
            CONTOH MERCHANT ID
          </div>
          <div className="flex flex-wrap gap-2">
            {merchants.slice(0, 8).map((m) => {
              const score = m.risk_score || 0;
              const color = riskColor(score);
              return (
                <span
                  key={m.id}
                  className="font-mono text-[11px] font-semibold rounded-sm px-2 py-0.5"
                  style={{
                    color,
                    background: `${color}15`,
                    border: `1px solid ${color}30`,
                  }}
                >
                  {m.id}
                </span>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}

import { useCallback } from 'react';
import RiskBadge, { riskColor } from './RiskBadge';
import { submitReport } from '../lib/api';

export default function MerchantDetail({ merchant, onClose }) {
  const handleReport = useCallback(async () => {
    if (!confirm(`Laporkan ${merchant.name} (${merchant.id}) ke OJK?`)) return;
    try {
      await submitReport(merchant.id, 'Auto-report from dashboard');
      alert('Laporan terkirim ke OJK.');
    } catch (e) {
      alert('Gagal mengirim laporan: ' + e.message);
    }
  }, [merchant]);

  if (!merchant) {
    return (
      <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg h-full flex flex-col items-center justify-center gap-3 opacity-40 p-5">
        <div className="w-10 h-10 rounded-full border border-dashed border-[rgba(255,255,255,0.25)] flex items-center justify-center text-lg text-[rgba(255,255,255,0.4)]">⊙</div>
        <div className="font-body text-xs text-[rgba(255,255,255,0.4)] text-center">Klik node pada graph untuk lihat detail merchant</div>
      </div>
    );
  }

  const score = merchant.risk_score || merchant.risk || 0;
  const color = riskColor(score);

  return (
    <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg p-5">
      <div className="flex justify-between items-start mb-4">
        <div>
          <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3" style={{ letterSpacing: '0.2em' }}>
            DETAIL
          </div>
          <div className="font-mono text-sm font-semibold text-white">{merchant.name}</div>
          <div className="font-mono text-[13px] text-[rgba(255,255,255,0.55)] mt-1">{merchant.id}</div>
        </div>
        <button onClick={onClose} className="font-mono text-xs text-[rgba(255,255,255,0.4)] hover:text-white cursor-pointer">&times;</button>
      </div>

      <div className="mb-5">
        <div className="flex justify-between items-center mb-1.5">
          <span className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase" style={{ letterSpacing: '0.15em' }}>Risk Score</span>
          <RiskBadge score={score} />
        </div>
        <div className="bg-[#07070F] rounded h-1.5 overflow-hidden">
          <div
            className="h-full rounded transition-all duration-500"
            style={{ width: `${score}%`, background: `linear-gradient(90deg, #30D158, ${color})` }}
          />
        </div>
        <div className="font-mono text-[22px] font-bold mt-2" style={{ color }}>
          {score}<span className="text-[13px] text-[rgba(255,255,255,0.4)]">/100</span>
        </div>
      </div>

      <div className="flex flex-col gap-2.5">
        {[
          ['MCC', merchant.mcc || '-'],
          ['NPWP', merchant.npwp ? 'VALID' : 'TIDAK TERDAFTAR'],
          ['LAPORAN', `${merchant.reports || 0}x`],
          ['STATUS', (merchant.flagged || score >= 70) ? 'TERFLAG' : 'AMAN'],
        ].map(([k, v]) => (
          <div key={k} className="flex justify-between pb-2.5 border-b border-[rgba(255,255,255,0.04)]">
            <span className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase" style={{ letterSpacing: '0.15em' }}>{k}</span>
            <span className="font-mono text-[13px]" style={{ color: k === 'NPWP' && !merchant.npwp ? '#FF2D55' : '#FFFFFF' }}>{v}</span>
          </div>
        ))}
      </div>

      {(score >= 70 || merchant.flagged) && (
        <button
          onClick={handleReport}
          className="mt-4 w-full font-mono text-[11px] font-bold uppercase tracking-wider py-2.5 rounded-md cursor-pointer transition-all"
          style={{
            background: 'rgba(255,45,85,0.12)',
            color: '#FF2D55',
            border: '1px solid rgba(255,45,85,0.3)',
          }}
        >
          ↗ LAPORKAN KE OJK
        </button>
      )}
    </div>
  );
}

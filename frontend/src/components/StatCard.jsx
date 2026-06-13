export default function StatCard({ label, value, accent }) {
  return (
    <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg p-5 flex-1 min-w-[130px]">
      <div
        className="font-mono text-[36px] font-bold leading-none"
        style={{ color: accent || '#2B2BFF' }}
      >
        {value.toLocaleString()}
      </div>
      <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] mt-1.5 uppercase"
        style={{ letterSpacing: '0.2em' }}
      >
        {label}
      </div>
    </div>
  );
}

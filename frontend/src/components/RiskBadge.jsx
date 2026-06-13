export function riskColor(score) {
  if (score >= 70) return '#FF2D55';
  if (score >= 30) return '#FF9F0A';
  return '#30D158';
}

export function riskLabel(score) {
  if (score >= 70) return 'MERAH';
  if (score >= 30) return 'KUNING';
  return 'HIJAU';
}

export default function RiskBadge({ score }) {
  const label = riskLabel(score);
  const color = riskColor(score);
  const bgMap = { MERAH: 'rgba(255,45,85,0.12)', KUNING: 'rgba(255,159,10,0.12)', HIJAU: 'rgba(48,209,88,0.12)' };
  return (
    <span
      className="font-mono text-[10px] font-bold px-2 py-[3px] rounded-sm uppercase"
      style={{
        color,
        background: bgMap[label],
        border: `1px solid ${color}40`,
        letterSpacing: '0.15em',
      }}
    >
      {label}
    </span>
  );
}

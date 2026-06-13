import { useRef, useEffect, useCallback } from 'react';
import { riskColor } from './RiskBadge';

export default function NetworkGraph({ merchants, selected, onSelect }) {
  const canvasRef = useRef(null);
  const animRef = useRef(null);
  const nodesRef = useRef([]);
  const timeRef = useRef(0);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');

    const resize = () => {
      const dpr = window.devicePixelRatio || 1;
      canvas.width = canvas.offsetWidth * dpr;
      canvas.height = canvas.offsetHeight * dpr;
      ctx.scale(dpr, dpr);
    };
    resize();
    window.addEventListener('resize', resize);

    nodesRef.current = merchants.map((m) => ({
      ...m,
      px: (m.x || Math.random()) * canvas.offsetWidth,
      py: (m.y || Math.random()) * canvas.offsetHeight,
      vx: (Math.random() - 0.5) * 0.3,
      vy: (Math.random() - 0.5) * 0.3,
      pulse: Math.random() * Math.PI * 2,
    }));

    const draw = () => {
      timeRef.current += 0.02;
      const w = canvas.offsetWidth;
      const h = canvas.offsetHeight;
      ctx.clearRect(0, 0, w, h);

      // Background grid — blue-tinted
      ctx.strokeStyle = 'rgba(43,43,255,0.04)';
      ctx.lineWidth = 1;
      for (let x = 0; x < w; x += 40) {
        ctx.beginPath(); ctx.moveTo(x, 0); ctx.lineTo(x, h); ctx.stroke();
      }
      for (let y = 0; y < h; y += 40) {
        ctx.beginPath(); ctx.moveTo(0, y); ctx.lineTo(w, y); ctx.stroke();
      }

      const nodes = nodesRef.current;

      // Drift
      nodes.forEach((n) => {
        n.px += n.vx;
        n.py += n.vy;
        if (n.px < 30 || n.px > w - 30) n.vx *= -1;
        if (n.py < 30 || n.py > h - 30) n.vy *= -1;
      });

      // Edges
      const risky = nodes.filter((n) => (n.risk_score || n.risk || 0) >= 30 || n.flagged);
      risky.forEach((a, i) => {
        risky.forEach((b, j) => {
          if (j <= i) return;
          const dist = Math.hypot(a.px - b.px, a.py - b.py);
          if (dist > 180) return;
          const alpha = (1 - dist / 180) * 0.35;
          ctx.beginPath();
          ctx.strokeStyle = `rgba(255,45,85,${alpha})`;
          ctx.lineWidth = 1;
          ctx.setLineDash([4, 4]);
          ctx.lineDashOffset = -timeRef.current * 8;
          ctx.moveTo(a.px, a.py);
          ctx.lineTo(b.px, b.py);
          ctx.stroke();
          ctx.setLineDash([]);
        });
      });

      // Nodes
      nodes.forEach((n) => {
        const score = n.risk_score || n.risk || 0;
        const isSelected = selected === n.id;
        const isFlagged = n.flagged || score >= 70;
        const r = isFlagged ? 10 : 7;
        const pulse = Math.sin(timeRef.current * 2 + n.pulse);
        const color = riskColor(score);

        // Pulse ring
        if (isFlagged) {
          ctx.beginPath();
          ctx.arc(n.px, n.py, r + 6 + pulse * 4, 0, Math.PI * 2);
          ctx.strokeStyle = `rgba(255,45,85,${0.15 + pulse * 0.08})`;
          ctx.lineWidth = 1.5;
          ctx.stroke();
        }

        // Glow
        const glow = ctx.createRadialGradient(n.px, n.py, 0, n.px, n.py, r * 3);
        glow.addColorStop(0, `${color}44`);
        glow.addColorStop(1, 'transparent');
        ctx.beginPath();
        ctx.arc(n.px, n.py, r * 3, 0, Math.PI * 2);
        ctx.fillStyle = glow;
        ctx.fill();

        // Core
        ctx.beginPath();
        ctx.arc(n.px, n.py, isSelected ? r + 3 : r, 0, Math.PI * 2);
        ctx.fillStyle = color;
        ctx.fill();

        if (isSelected) {
          ctx.beginPath();
          ctx.arc(n.px, n.py, r + 6, 0, Math.PI * 2);
          ctx.strokeStyle = '#2B2BFF';
          ctx.lineWidth = 2;
          ctx.stroke();
        }

        // Label
        ctx.font = "10px 'JetBrains Mono', monospace";
        ctx.fillStyle = 'rgba(255,255,255,0.45)';
        ctx.fillText(n.id || n.name, n.px + r + 4, n.py + 4);
      });

      animRef.current = requestAnimationFrame(draw);
    };

    draw();

    return () => {
      cancelAnimationFrame(animRef.current);
      window.removeEventListener('resize', resize);
    };
  }, [merchants]);

  const handleClick = useCallback((e) => {
    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const mx = e.clientX - rect.left;
    const my = e.clientY - rect.top;
    const nodes = nodesRef.current;
    for (const n of nodes) {
      if (Math.hypot(n.px - mx, n.py - my) < 16) {
        onSelect(n.id === selected ? null : n.id);
        return;
      }
    }
    onSelect(null);
  }, [selected, onSelect]);

  return (
    <canvas
      ref={canvasRef}
      onClick={handleClick}
      className="w-full h-full block"
      style={{ cursor: 'crosshair' }}
    />
  );
}

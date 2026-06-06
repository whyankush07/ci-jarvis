import React from "react";

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export const Card = ({ children, className = "" }: CardProps) => {
  return (
    <div className={`bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden ${className}`}>
      {children}
    </div>
  );
};

export const CardHeader = ({ children, className = "" }: CardProps) => (
  <div className={`px-6 py-4 border-b border-zinc-100 dark:border-zinc-800 ${className}`}>
    {children}
  </div>
);

export const CardContent = ({ children, className = "" }: CardProps) => (
  <div className={`px-6 py-4 ${className}`}>
    {children}
  </div>
);

export const CardFooter = ({ children, className = "" }: CardProps) => (
  <div className={`px-6 py-4 bg-zinc-50/50 dark:bg-zinc-800/50 border-t border-zinc-100 dark:border-zinc-800 ${className}`}>
    {children}
  </div>
);

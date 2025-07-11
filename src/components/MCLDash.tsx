export const MCLDash = () => {
    return (
        <div>
            <div className="flex justify-between gap-2 w-full">
                <p>logo</p>
                <div>
                    <h4 className="text-orange">MAP</h4>
                    <h4>2</h4>
                </div>
                <div>
                    <h4 className="text-orange">TC1</h4>
                    <h4>2</h4>
                </div>
                <div>
                    <h4 className="text-orange">TC2</h4>
                    <h4>2</h4>
                </div>
                <div>
                    <h4 className="text-orange">ABS</h4>
                    <h4>0</h4>
                </div>
                <div>
                    <h4 className="text-orange">BB</h4>
                    <h4 className="text-orange">56.8</h4>
                </div>
                <p>logo2</p>
            </div>
            <div className="grid grid-cols-3">
                <div>
                    <McLarenSwoosh />
                    <div className="px-4">
                        <h3 className="text-orange">PREDICTED LAPTIME</h3>
                        <h2>01:28:36</h2>
                        <h3 className="text-orange">LAP DELTA</h3>
                        <h2 className="text-green">-00:00</h2>
                        <div className="flex justify-between">
                            <div>
                                <h3 className="text-orange">LAST LAP</h3>
                                <h3>01:29:88</h3>
                            </div>
                            <div>
                                <h3 className="text-orange">LAPS</h3>
                                <h3>2</h3>
                            </div>
                        </div>
                        <div className="flex justify-between">
                            <div>
                                <h3 className="text-orange">WATER TEMP</h3>
                                <h3>83</h3>
                            </div>
                            <div>
                                <h3 className="text-orange">BATTERY VOLTS</h3>
                                <h3>12.1</h3>
                            </div>
                        </div>

                    </div>
                    <McLarenSwoosh />
                </div>
                <div>
                    <h3 className="text-orange">RPM</h3>
                    <h3>8234</h3>
                    <h1 className="text-orange !text-[140px]">6</h1>
                    <h3>268</h3>
                    <h3 className="text-orange">SPEED</h3>
                </div>
                <div>
                    <McLarenSwoosh />
                    <div className="px-4">
                        <h3 className="text-orange">TYRE DATA</h3>
                        <div>
                            <div className="grid grid-cols-4 p-0.5 justify-items-center">
                                <div>
                                    <h4>1.34</h4>
                                    <h4>84</h4>
                                </div>
                                <div className="h-[34px] w-[24px] rounded-sm bg-blue-300" />
                                <div className="h-[34px] w-[24px] rounded-sm bg-blue-300" />
                                <div>
                                    <h4>1.34</h4>
                                    <h4>84</h4>
                                </div>
                            </div>
                            <div className="grid grid-cols-4 p-0.5 justify-items-center">
                                <div>
                                    <h4>1.34</h4>
                                    <h4>84</h4>
                                </div>
                                <div className="h-[34px] w-[24px] rounded-sm bg-blue-300" />
                                <div className="h-[34px] w-[24px] rounded-sm bg-blue-300" />
                                <div>
                                    <h4>1.34</h4>
                                    <h4>84</h4>
                                </div>
                            </div>
                        </div>
                        <div className="flex justify-between">
                            <div>
                                <h3 className="text-orange">FUEL USED</h3>
                                <h3>5.5</h3>
                            </div>
                            <div>
                                <h3 className="text-orange">FUEL LEVEL</h3>
                                <h3 className="text-orange">26.4</h3>
                            </div>
                        </div>
                        <h3 className="text-orange">FUEL LAST LAP</h3>
                        <h3>2.416</h3>

                    </div>
                    <McLarenSwoosh />
                </div>
            </div>
        </div>
    )
}

const McLarenSwoosh = () => {

    return (
        <svg width="200" height="25" viewBox="0 0 400 25">
            <defs>
                <linearGradient id="mclarenLine" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" stopColor="transparent" />
                    <stop offset="10%" stopColor="#FF4400" stopOpacity="1" />
                    <stop offset="30%" stopColor="#FF6600" stopOpacity="1" />
                    <stop offset="50%" stopColor="#FFAA00" stopOpacity="1" />
                    <stop offset="70%" stopColor="#FF6600" stopOpacity="1" />
                    <stop offset="90%" stopColor="#FF4400" stopOpacity="1" />
                    <stop offset="100%" stopColor="transparent" />
                </linearGradient>
                <filter id="glow">
                    <feGaussianBlur stdDeviation="2" result="coloredBlur" />
                    <feMerge>
                        <feMergeNode in="coloredBlur" />
                        <feMergeNode in="SourceGraphic" />
                    </feMerge>
                </filter>
            </defs>
            <path
                d="M 3 16.7 L 39.4 4.5 H 241 H 383 H 400"
                stroke="url(#mclarenLine)"
                strokeWidth="2.5"
                fill="none"
                strokeLinecap="round"
                filter="url(#glow)"
            />
        </svg>
    );
};

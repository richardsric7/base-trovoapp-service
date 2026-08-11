document.addEventListener('DOMContentLoaded', () => {
    // Hero Background Switching
    const heroCards = document.querySelectorAll('.hero-card');
    const heroBgs = document.querySelectorAll('.hero-bg');
    const heroSection = document.querySelector('.hero');

    heroCards.forEach(card => {
        const handleInteraction = () => {
            const index = card.getAttribute('data-index');
            
            // Update active card
            heroCards.forEach(c => c.classList.remove('active'));
            card.classList.add('active');

            // Update active background
            heroBgs.forEach(bg => bg.classList.remove('active'));
            const targetBg = document.getElementById(`hero-bg-${index}`);
            if (targetBg) targetBg.classList.add('active');
        };

        card.addEventListener('mouseenter', handleInteraction);
        card.addEventListener('click', handleInteraction);
    });

    // Step Items Scroll Activation
    const stepItems = document.querySelectorAll('.step-item');
    const stepsSection = document.querySelector('.steps-section');

    function setActiveStepByScroll() {
        if (!stepItems.length || !stepsSection) return;

        const sectionRect = stepsSection.getBoundingClientRect();
        const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
        const sectionInView = sectionRect.top < viewportHeight && sectionRect.bottom > 0;

        if (!sectionInView) return;

        // Use a focus line inside the viewport to decide which step is "current"
        const focusY = viewportHeight * 0.45;
        let activeIndex = 0;
        let closestDistance = Infinity;

        stepItems.forEach((item, index) => {
            const rect = item.getBoundingClientRect();
            const itemCenter = rect.top + rect.height / 2;
            const distance = Math.abs(itemCenter - focusY);

            if (distance < closestDistance) {
                closestDistance = distance;
                activeIndex = index;
            }
        });

        stepItems.forEach((item, index) => {
            item.classList.toggle('active', index === activeIndex);
        });
    }

    window.addEventListener('scroll', setActiveStepByScroll, { passive: true });
    window.addEventListener('resize', setActiveStepByScroll);
    setActiveStepByScroll();

    // Market Section Tabs + Auto Progress
    const marketTabs = Array.from(document.querySelectorAll('.market-tab'));
    const marketContents = Array.from(document.querySelectorAll('.market-content'));
    const marketFeatureDuration = 2400;
    let marketActiveIndex = marketTabs.findIndex(tab => tab.classList.contains('active'));
    let marketActiveFeatureIndex = 0;
    let marketProgressStart = null;
    let marketRafId = null;

    if (marketActiveIndex < 0) marketActiveIndex = 0;

    function getActiveMarketFeatures() {
        const activeContent = marketContents[marketActiveIndex];
        if (!activeContent) return [];
        return Array.from(activeContent.querySelectorAll('.market-feature'));
    }

    function resetAllMarketFeatureProgress() {
        marketContents.forEach(content => {
            const features = content.querySelectorAll('.market-feature');
            features.forEach(feature => {
                feature.classList.remove('progress-active');
                feature.style.setProperty('--feature-progress', '0%');
            });
        });
    }

    function updateMarketImageForFeature(content, feature) {
        if (!content || !feature) return;
        const image = content.querySelector('.market-image img');
        if (!image) return;

        const featureImage = feature.getAttribute('data-feature-image');
        const featureAlt = feature.getAttribute('data-feature-alt');

        if (featureImage) image.src = featureImage;
        if (featureAlt) image.alt = featureAlt;
    }

    function alignMarketImageToFeatureStack(content) {
        if (!content) return;
        const features = Array.from(content.querySelectorAll('.market-feature'));
        const imageWrap = content.querySelector('.market-image');
        if (!imageWrap || features.length === 0) return;

        const firstRect = features[0].getBoundingClientRect();
        const lastRect = features[features.length - 1].getBoundingClientRect();
        const contentRect = content.getBoundingClientRect();

        const topOffset = Math.max(0, firstRect.top - contentRect.top);
        const stackHeight = Math.max(0, lastRect.bottom - firstRect.top);

        imageWrap.style.setProperty('--market-image-top', `${topOffset}px`);
        imageWrap.style.setProperty('--market-image-height', `${stackHeight}px`);

        // Fine-tune after layout so image top exactly matches first feature top.
        requestAnimationFrame(() => {
            const refreshedFirstRect = features[0].getBoundingClientRect();
            const imageRect = imageWrap.getBoundingClientRect();
            const delta = imageRect.top - refreshedFirstRect.top;
            if (Math.abs(delta) > 0.5) {
                const correctedTop = Math.max(0, topOffset - delta);
                imageWrap.style.setProperty('--market-image-top', `${correctedTop}px`);
            }
        });
    }

    function setMarketFeatureProgress(features, featureIndex, progressPercent) {
        const clamped = Math.max(0, Math.min(100, progressPercent));
        const safeProgress = `${clamped}%`;
        const safeScale = `${clamped / 100}`;
        const activeContent = marketContents[marketActiveIndex];

        features.forEach((feature, index) => {
            const isActive = index === featureIndex;
            feature.classList.toggle('progress-active', isActive);
            feature.style.setProperty('--feature-progress', isActive ? safeProgress : '0%');
            feature.style.setProperty('--feature-progress-scale', isActive ? safeScale : '0');
        });

        if (features[featureIndex] && activeContent) {
            updateMarketImageForFeature(activeContent, features[featureIndex]);
        }
    }

    function activateMarketByIndex(index, restartProgress = true) {
        if (!marketTabs.length || !marketContents.length) return;
        const safeIndex = ((index % marketTabs.length) + marketTabs.length) % marketTabs.length;
        marketActiveIndex = safeIndex;
        const target = marketTabs[safeIndex].getAttribute('data-target');

        marketTabs.forEach((tab, i) => {
            tab.classList.toggle('active', i === safeIndex);
        });

        marketContents.forEach(content => {
            content.classList.toggle('active', content.id === target);
        });

        if (restartProgress) {
            resetAllMarketFeatureProgress();
            marketActiveFeatureIndex = 0;
            marketProgressStart = performance.now();
            const features = getActiveMarketFeatures();
            if (features.length) {
                setMarketFeatureProgress(features, marketActiveFeatureIndex, 0);
            }
        }

        requestAnimationFrame(() => {
            alignMarketImageToFeatureStack(marketContents[safeIndex]);
        });
    }

    function runMarketProgress(now) {
        if (!marketTabs.length) return;
        if (marketProgressStart === null) marketProgressStart = now;

        const activeFeatures = getActiveMarketFeatures();
        if (!activeFeatures.length) {
            activateMarketByIndex(marketActiveIndex + 1, true);
            marketRafId = requestAnimationFrame(runMarketProgress);
            return;
        }

        const elapsed = now - marketProgressStart;
        const progress = (elapsed / marketFeatureDuration) * 100;
        setMarketFeatureProgress(activeFeatures, marketActiveFeatureIndex, progress);

        if (elapsed >= marketFeatureDuration) {
            if (marketActiveFeatureIndex < activeFeatures.length - 1) {
                marketActiveFeatureIndex += 1;
                marketProgressStart = now;
                setMarketFeatureProgress(activeFeatures, marketActiveFeatureIndex, 0);
            } else {
                activateMarketByIndex(marketActiveIndex + 1, true);
            }
        }

        marketRafId = requestAnimationFrame(runMarketProgress);
    }

    marketTabs.forEach((tab, index) => {
        tab.addEventListener('click', () => {
            activateMarketByIndex(index, true);
        });
    });

    if (marketTabs.length && marketContents.length) {
        activateMarketByIndex(marketActiveIndex, true);
        marketRafId = requestAnimationFrame(runMarketProgress);
        window.addEventListener('resize', () => {
            alignMarketImageToFeatureStack(marketContents[marketActiveIndex]);
        });
    }

    // Compliance Cards Toggle
    const complianceCards = document.querySelectorAll('.compliance-card');
    const complianceButtons = document.querySelectorAll('.compliance-toggle');

    function setComplianceState(activeCard) {
        complianceCards.forEach((card) => {
            const button = card.querySelector('.compliance-toggle');
            const isActive = card === activeCard;
            card.classList.toggle('active', isActive);
            if (button) {
                button.innerHTML = isActive
                    ? '<i class="fas fa-times"></i> Close'
                    : 'Learn More';
            }
        });
    }

    complianceButtons.forEach((button) => {
        button.addEventListener('click', () => {
            const card = button.closest('.compliance-card');
            if (!card) return;

            if (card.classList.contains('active')) {
                setComplianceState(null);
                return;
            }

            setComplianceState(card);
        });
    });

    // Scroll Header Effect
    const header = document.querySelector('header');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header.style.backgroundColor = 'rgba(255, 255, 255, 0.95)';
            header.style.padding = '10px 0';
            header.style.boxShadow = '0 4px 20px rgba(0,0,0,0.05)';
        } else {
            header.style.backgroundColor = 'rgba(255, 255, 255, 0.9)';
            header.style.padding = '20px 0';
            header.style.boxShadow = 'none';
        }
    });

    // Simple reveal animation on scroll
    // Keep step items visible at all times so the "How it works" copy
    // does not disappear if IntersectionObserver timing varies.
    const revealElements = document.querySelectorAll('.compliance-card, .market-content');
    const revealObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
                revealObserver.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1 });

    revealElements.forEach(el => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'all 0.6s ease-out';
        revealObserver.observe(el);
    });
});

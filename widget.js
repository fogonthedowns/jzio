/**
 * HelpPet Booking Widget
 * Embeddable appointment booking widget for third-party websites
 * Version: 1.0.0
 */

(function() {
  'use strict';

  const API_BASE_URL = 'http://localhost:8000';  // Change to https://api.helppet.ai for production

  class HelpPetWidget {
    constructor(config) {
      this.config = {
        practiceId: config.practiceId,
        token: config.token,
        containerId: config.containerId,
        primaryColor: config.primaryColor || '#0088cc',
        onBookingComplete: config.onBookingComplete || (() => {}),
        onBookingError: config.onBookingError || (() => {})
      };

      this.currentStep = 1;
      this.formData = {};
      this.container = null;
      this.shadowRoot = null;
      this.currentDate = new Date(); // Track the date being viewed
      
      this.init();
    }

    init() {
      console.log('HelpPet Widget: Initializing...', {
        practiceId: this.config.practiceId,
        token: this.config.token.substring(0, 10) + '...',
        containerId: this.config.containerId
      });

      // Find container
      const container = document.getElementById(this.config.containerId);
      if (!container) {
        console.error(`HelpPet Widget: Container #${this.config.containerId} not found`);
        return;
      }

      // Create shadow DOM for style isolation
      this.container = container;
      this.shadowRoot = container.attachShadow({ mode: 'open' });

      // Inject styles
      this.injectStyles();

      // Render initial step
      this.render();
      console.log('HelpPet Widget: Initialized successfully on Step 1');
    }

    injectStyles() {
      const style = document.createElement('style');
      style.textContent = `
        * {
          box-sizing: border-box;
          margin: 0;
          padding: 0;
        }

        .widget-container {
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
          max-width: 600px;
          margin: 0 auto;
          padding: 24px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
        }

        .widget-header {
          text-align: center;
          margin-bottom: 32px;
        }

        .widget-title {
          font-size: 24px;
          font-weight: 600;
          color: #1a1a1a;
          margin-bottom: 8px;
        }

        .widget-subtitle {
          font-size: 14px;
          color: #666;
        }

        .step-indicator {
          display: flex;
          justify-content: center;
          gap: 12px;
          margin-bottom: 32px;
        }

        .step-dot {
          width: 8px;
          height: 8px;
          border-radius: 50%;
          background: #e0e0e0;
          transition: all 0.3s;
        }

        .step-dot.active {
          background: ${this.config.primaryColor};
          width: 24px;
          border-radius: 4px;
        }

        .form-group {
          margin-bottom: 20px;
        }

        .form-label {
          display: block;
          font-size: 14px;
          font-weight: 500;
          color: #333;
          margin-bottom: 8px;
        }

        .form-input,
        .form-select {
          width: 100%;
          padding: 12px;
          font-size: 14px;
          border: 1px solid #ddd;
          border-radius: 8px;
          transition: border-color 0.2s;
        }

        .form-input:focus,
        .form-select:focus {
          outline: none;
          border-color: ${this.config.primaryColor};
        }

        .form-error {
          color: #dc3545;
          font-size: 12px;
          margin-top: 4px;
        }

        .button-group {
          display: flex;
          gap: 12px;
          margin-top: 32px;
        }

        .btn {
          flex: 1;
          padding: 14px 24px;
          font-size: 16px;
          font-weight: 500;
          border: none;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.2s;
        }

        .btn-primary {
          background: ${this.config.primaryColor};
          color: white;
        }

        .btn-primary:hover {
          opacity: 0.9;
        }

        .btn-primary:disabled {
          background: #ccc;
          cursor: not-allowed;
        }

        .btn-secondary {
          background: #f5f5f5;
          color: #333;
        }

        .btn-secondary:hover {
          background: #e0e0e0;
        }

        .calendar-container {
          margin: 20px 0;
        }

        .time-slots {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
          gap: 8px;
          margin-top: 16px;
        }

        .time-slot {
          padding: 12px;
          text-align: center;
          border: 1px solid #ddd;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.2s;
        }

        .time-slot:hover {
          border-color: ${this.config.primaryColor};
          background: ${this.config.primaryColor}10;
        }

        .time-slot.selected {
          background: ${this.config.primaryColor};
          color: white;
          border-color: ${this.config.primaryColor};
        }

        .time-slot.disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }

        .confirmation {
          text-align: center;
          padding: 40px 20px;
        }

        .success-icon {
          width: 64px;
          height: 64px;
          background: #28a745;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 24px;
          color: white;
          font-size: 32px;
        }

        .confirmation-code {
          font-size: 24px;
          font-weight: 600;
          color: ${this.config.primaryColor};
          margin: 16px 0;
        }

        .loading {
          text-align: center;
          padding: 40px;
          color: #666;
        }

        .error-message {
          background: #fff3cd;
          border: 1px solid #ffc107;
          border-radius: 8px;
          padding: 16px;
          margin-bottom: 20px;
        }

        @media (max-width: 640px) {
          .widget-container {
            padding: 16px;
          }

          .time-slots {
            grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
          }
        }
      `;
      // Store reference to style element
      this.styleElement = style;
      this.shadowRoot.appendChild(style);
    }

    render() {
      const container = document.createElement('div');
      container.className = 'widget-container';
      
      container.innerHTML = `
        <div class="widget-header">
          <h2 class="widget-title">Book an Appointment</h2>
          <p class="widget-subtitle">Step ${this.currentStep} of 4</p>
        </div>
        <div class="step-indicator">
          <div class="step-dot ${this.currentStep >= 1 ? 'active' : ''}"></div>
          <div class="step-dot ${this.currentStep >= 2 ? 'active' : ''}"></div>
          <div class="step-dot ${this.currentStep >= 3 ? 'active' : ''}"></div>
          <div class="step-dot ${this.currentStep >= 4 ? 'active' : ''}"></div>
        </div>
        <div class="widget-content">
          ${this.renderStep()}
        </div>
      `;

      // Clear and append (keep style element)
      this.shadowRoot.innerHTML = '';
      this.shadowRoot.appendChild(this.styleElement);
      this.shadowRoot.appendChild(container);

      // Attach event listeners
      this.attachEventListeners();
    }

    renderStep() {
      switch (this.currentStep) {
        case 1:
          return this.renderStep1();
        case 2:
          return this.renderStep2();
        case 3:
          return this.renderStep3();
        case 4:
          return this.renderStep4();
        default:
          return '';
      }
    }

    renderStep1() {
      return `
        <div class="form-group">
          <label class="form-label">Are you a new or returning client?</label>
          <select class="form-select" id="client-status">
            <option value="">Select...</option>
            <option value="new" ${this.formData.clientStatus === 'new' ? 'selected' : ''}>New Client</option>
            <option value="returning" ${this.formData.clientStatus === 'returning' ? 'selected' : ''}>Returning Client</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">Pet Species</label>
          <select class="form-select" id="pet-species">
            <option value="">Select...</option>
            <option value="dog" ${this.formData.species === 'dog' ? 'selected' : ''}>Dog</option>
            <option value="cat" ${this.formData.species === 'cat' ? 'selected' : ''}>Cat</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">Appointment Type</label>
          <select class="form-select" id="appointment-type">
            <option value="">Select...</option>
            <option value="checkup" ${this.formData.appointmentType === 'checkup' ? 'selected' : ''}>Wellness Checkup</option>
            <option value="emergency" ${this.formData.appointmentType === 'emergency' ? 'selected' : ''}>Emergency</option>
            <option value="surgery" ${this.formData.appointmentType === 'surgery' ? 'selected' : ''}>Surgery</option>
            <option value="consultation" ${this.formData.appointmentType === 'consultation' ? 'selected' : ''}>Consultation</option>
          </select>
        </div>

        <div class="button-group">
          <button class="btn btn-primary" id="next-btn">Next</button>
        </div>
      `;
    }

    renderStep2() {
      return `
        <div class="loading">Loading available times...</div>
        <div class="button-group">
          <button class="btn btn-secondary" id="back-btn">Back</button>
          <button class="btn btn-primary" id="next-btn" disabled>Select a Time</button>
        </div>
      `;
    }

    renderStep3() {
      return `
        <div class="form-group">
          <label class="form-label">First Name</label>
          <input type="text" class="form-input" id="first-name" value="${this.formData.firstName || ''}" />
        </div>

        <div class="form-group">
          <label class="form-label">Last Name</label>
          <input type="text" class="form-input" id="last-name" value="${this.formData.lastName || ''}" />
        </div>

        <div class="form-group">
          <label class="form-label">Email</label>
          <input type="email" class="form-input" id="email" value="${this.formData.email || ''}" />
        </div>

        <div class="form-group">
          <label class="form-label">Phone Number</label>
          <input type="tel" class="form-input" id="phone" value="${this.formData.phone || ''}" placeholder="(555) 123-4567" />
        </div>

        <div class="form-group">
          <label class="form-label">Pet Name</label>
          <input type="text" class="form-input" id="pet-name" value="${this.formData.petName || ''}" />
        </div>

        <div class="form-group">
          <label class="form-label">Reason for Visit (Optional)</label>
          <input type="text" class="form-input" id="reason" value="${this.formData.reason || ''}" placeholder="Annual checkup, vaccination, etc." />
        </div>

        <div class="button-group">
          <button class="btn btn-secondary" id="back-btn">Back</button>
          <button class="btn btn-primary" id="submit-btn">Book Appointment</button>
        </div>
      `;
    }

    renderStep4() {
      return `
        <div class="confirmation">
          <div class="success-icon">✓</div>
          <h3 class="widget-title">Appointment Confirmed!</h3>
          <div class="confirmation-code">${this.formData.confirmationCode || 'PENDING'}</div>
          <p style="color: #666; margin-top: 16px;">
            A confirmation email has been sent to<br/>
            <strong>${this.formData.email}</strong>
          </p>
          <p style="color: #666; margin-top: 24px; font-size: 14px;">
            <strong>${this.formData.displayTime}</strong>
          </p>
          <p style="color: #999; margin-top: 8px; font-size: 12px;">
            (Times shown in your local timezone)
          </p>
        </div>
      `;
    }

    attachEventListeners() {
      const nextBtn = this.shadowRoot.getElementById('next-btn');
      const backBtn = this.shadowRoot.getElementById('back-btn');
      const submitBtn = this.shadowRoot.getElementById('submit-btn');

      if (nextBtn) {
        nextBtn.addEventListener('click', () => this.handleNext());
      }

      if (backBtn) {
        backBtn.addEventListener('click', () => this.handleBack());
      }

      if (submitBtn) {
        submitBtn.addEventListener('click', () => this.handleSubmit());
      }

      // Step-specific listeners
      if (this.currentStep === 1) {
        const inputs = ['client-status', 'pet-species', 'appointment-type'];
        inputs.forEach(id => {
          const input = this.shadowRoot.getElementById(id);
          if (input) {
            input.addEventListener('change', () => this.validateStep1());
          }
        });
      }
    }

    validateStep1() {
      const clientStatus = this.shadowRoot.getElementById('client-status').value;
      const species = this.shadowRoot.getElementById('pet-species').value;
      const appointmentType = this.shadowRoot.getElementById('appointment-type').value;
      
      const nextBtn = this.shadowRoot.getElementById('next-btn');
      nextBtn.disabled = !(clientStatus && species && appointmentType);
    }

    handleNext() {
      if (this.currentStep === 1) {
        // Save step 1 data
        this.formData.clientStatus = this.shadowRoot.getElementById('client-status').value;
        this.formData.species = this.shadowRoot.getElementById('pet-species').value;
        this.formData.appointmentType = this.shadowRoot.getElementById('appointment-type').value;
      }

      this.currentStep++;
      this.render();

      if (this.currentStep === 2) {
        this.loadAvailability();
      }
    }

    handleBack() {
      // Guard against going below step 1
      if (this.currentStep <= 1) return;
      
      this.currentStep--;
      this.render();
    }

    async handleSubmit() {
      // Collect and trim form data
      this.formData.firstName = this.shadowRoot.getElementById('first-name').value.trim();
      this.formData.lastName = this.shadowRoot.getElementById('last-name').value.trim();
      this.formData.email = this.shadowRoot.getElementById('email').value.trim();
      this.formData.phone = this.shadowRoot.getElementById('phone').value.trim();
      this.formData.petName = this.shadowRoot.getElementById('pet-name').value.trim();
      this.formData.reason = this.shadowRoot.getElementById('reason').value.trim();

      // Basic validation
      if (!this.formData.firstName || !this.formData.lastName || !this.formData.email || !this.formData.phone || !this.formData.petName) {
        this.showInlineError('Please fill out all required fields.');
        return;
      }

      // Show loading state
      const submitBtn = this.shadowRoot.getElementById('submit-btn');
      const originalText = submitBtn.textContent;
      submitBtn.disabled = true;
      submitBtn.textContent = 'Booking...';

      // Submit booking
      try {
        const url = `${API_BASE_URL}/api/v1/public/appointments/book`;
        
        // Use the date/time already in practice timezone (stored from slot selection)
        const payload = {
          practice_id: this.config.practiceId,
          customer_name: `${this.formData.firstName} ${this.formData.lastName}`,
          customer_email: this.formData.email,
          customer_phone: this.formData.phone,
          pet_name: this.formData.petName,
          pet_species: this.formData.species || 'dog',
          pet_breed: null,
          requested_date: this.formData.selectedDate, // Already in practice timezone
          requested_time: this.formData.selectedTime, // Already in practice timezone
          clinician_id: this.formData.clinicianId || null,
          appointment_type: this.formData.appointmentType || 'checkup',
          notes: this.formData.reason || ''
        };

        console.log('HelpPet Widget: Submitting booking to', url);
        console.log('HelpPet Widget: Payload', payload);
        console.log('HelpPet Widget: Token', this.config.token.substring(0, 10) + '...');

        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.config.token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(payload)
        });

        console.log('HelpPet Widget: Response status', response.status);

        const data = await response.json();
        console.log('HelpPet Widget: Response data', data);

        if (data.success) {
          this.formData.confirmationCode = data.confirmation_number;
          // Keep the display time we already have
          this.currentStep = 4;
          this.render();
          this.config.onBookingComplete(data);
        } else {
          console.error('HelpPet Widget: Booking failed', data.error);
          this.config.onBookingError(data.error);
          this.showInlineError(data.error.user_message || 'Booking failed. Please try again.');
          submitBtn.disabled = false;
          submitBtn.textContent = originalText;
        }
      } catch (error) {
        console.error('HelpPet Widget: Network error', error);
        console.error('HelpPet Widget: Error details', {
          name: error.name,
          message: error.message,
          stack: error.stack
        });
        this.config.onBookingError(error);
        this.showInlineError('Network error. Please check your connection and try again.');
        submitBtn.disabled = false;
        submitBtn.textContent = originalText;
      }
    }

    showInlineError(message) {
      const content = this.shadowRoot.querySelector('.widget-content');
      let errorDiv = content.querySelector('.error-message');
      
      if (!errorDiv) {
        errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        content.insertBefore(errorDiv, content.firstChild);
      }
      
      errorDiv.innerHTML = `<strong>Error:</strong> ${message}`;
      
      // Auto-remove after 5 seconds
      setTimeout(() => {
        if (errorDiv && errorDiv.parentNode) {
          errorDiv.remove();
        }
      }, 5000);
    }

    changeDate(days) {
      this.currentDate.setDate(this.currentDate.getDate() + days);
      this.loadAvailability();
    }

    async loadAvailability() {
      console.log('HelpPet Widget: Loading availability...');
      
      // Show loading state
      const content = this.shadowRoot.querySelector('.widget-content');
      content.innerHTML = '<div class="loading">Loading available times...</div>';
      
      try {
        // Validate token format
        if (!this.config.token || !this.config.token.startsWith('hw_') || this.config.token.includes('•')) {
          throw new Error('Invalid API token. Please use your real token (starts with hw_), not the obfuscated version.');
        }

        // Use the current viewing date
        const dateStr = this.currentDate.toISOString().split('T')[0];
        
        const url = `${API_BASE_URL}/api/v1/public/availability/${this.config.practiceId}?date=${dateStr}`;
        
        console.log('HelpPet Widget: Fetching availability from', url);
        console.log('HelpPet Widget: Using token:', this.config.token.substring(0, 15) + '...');
        
        const response = await fetch(url, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${this.config.token}`,
            'Content-Type': 'application/json'
          }
        });

        console.log('HelpPet Widget: Availability response status', response.status);

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({ detail: response.statusText }));
          throw new Error(`HTTP ${response.status}: ${JSON.stringify(errorData)}`);
        }

        const data = await response.json();
        console.log('HelpPet Widget: Availability data', data);

        // Render availability slots
        this.renderAvailabilitySlots(data);

      } catch (error) {
        console.error('HelpPet Widget: Failed to load availability', error);
        
        // Show error message
        const content = this.shadowRoot.querySelector('.widget-content');
        content.innerHTML = `
          <div class="error-message">
            <strong>Unable to load available times</strong>
            <p style="margin-top: 8px;">${error.message}</p>
          </div>
          <div class="button-group">
            <button class="btn btn-secondary" id="back-btn">Back</button>
            <button class="btn btn-primary" id="retry-btn">Try Again</button>
          </div>
        `;

        const backBtn = content.querySelector('#back-btn');
        const retryBtn = content.querySelector('#retry-btn');
        
        backBtn.addEventListener('click', () => this.handleBack());
        retryBtn.addEventListener('click', () => this.loadAvailability());
      }
    }

    renderAvailabilitySlots(availabilityData) {
      console.log('HelpPet Widget: Rendering availability slots', availabilityData);
      
      const content = this.shadowRoot.querySelector('.widget-content');
      
      // Format the current date for display
      const displayDate = this.currentDate.toLocaleDateString('en-US', { 
        weekday: 'long', 
        month: 'long', 
        day: 'numeric',
        year: 'numeric'
      });
      
      // Check if we can go to previous day (don't allow past dates)
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const viewingDate = new Date(this.currentDate);
      viewingDate.setHours(0, 0, 0, 0);
      const canGoPrev = viewingDate > today;
      
      // Parse available slots from API response
      const slots = availabilityData.slots || availabilityData.available_slots || [];
      
      // Date navigation HTML
      const dateNavHtml = `
        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; padding: 12px; background: #f5f5f5; border-radius: 8px;">
          <button class="btn btn-secondary" id="prev-day-btn" ${!canGoPrev ? 'disabled' : ''} style="flex: 0 0 auto; padding: 8px 16px;">
            ← Previous
          </button>
          <div style="text-align: center; font-weight: 500;">
            ${displayDate}
          </div>
          <button class="btn btn-secondary" id="next-day-btn" style="flex: 0 0 auto; padding: 8px 16px;">
            Next →
          </button>
        </div>
      `;
      
      if (slots.length === 0) {
        content.innerHTML = `
          ${dateNavHtml}
          <div class="error-message">
            <strong>No available times</strong>
            <p style="margin-top: 8px;">There are no available appointment slots for this date. Use the buttons above to check other dates.</p>
          </div>
          <div class="button-group">
            <button class="btn btn-secondary" id="back-btn">Back</button>
          </div>
        `;
        
        const backBtn = content.querySelector('#back-btn');
        const prevBtn = content.querySelector('#prev-day-btn');
        const nextBtn = content.querySelector('#next-day-btn');
        
        backBtn.addEventListener('click', () => this.handleBack());
        if (prevBtn && !prevBtn.disabled) {
          prevBtn.addEventListener('click', () => this.changeDate(-1));
        }
        if (nextBtn) {
          nextBtn.addEventListener('click', () => this.changeDate(1));
        }
        return;
      }
      
      // Build time slot HTML
      const slotsHtml = slots.map(slot => {
        // Parse UTC datetime
        const utcDate = new Date(slot.start_time);
        
        // Convert to local timezone
        const localDate = new Date(utcDate.toLocaleString('en-US', { timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone }));
        
        // Format display (just show time, date is in the header)
        const timeStr = localDate.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true });
        
        // Store both UTC (for API) and practice-timezone date/time (for booking)
        return `<div class="time-slot" 
                     data-start="${slot.start_time}" 
                     data-end="${slot.end_time}"
                     data-date="${slot.date}"
                     data-time="${slot.time}"
                     data-clinician="${slot.clinician_id || ''}"
                     data-display="${displayDate} at ${timeStr}">
                  ${timeStr}
                </div>`;
      }).join('');
      
      content.innerHTML = `
        ${dateNavHtml}
        <div class="form-group">
          <label class="form-label">Select a Time</label>
          <p class="widget-subtitle" style="margin-bottom: 16px;">All times shown in your local timezone</p>
          <div class="time-slots">
            ${slotsHtml}
          </div>
        </div>
        <div class="button-group">
          <button class="btn btn-secondary" id="back-btn">Back</button>
          <button class="btn btn-primary" id="continue-btn" disabled>Continue</button>
        </div>
      `;

      // Attach time slot listeners
      const slotElements = content.querySelectorAll('.time-slot');
      const continueBtn = content.querySelector('#continue-btn');
      const backBtn = content.querySelector('#back-btn');
      const prevBtn = content.querySelector('#prev-day-btn');
      const nextDayBtn = content.querySelector('#next-day-btn');

      slotElements.forEach(slotEl => {
        slotEl.addEventListener('click', () => {
          slotElements.forEach(s => s.classList.remove('selected'));
          slotEl.classList.add('selected');
          
          // Store both UTC times and practice-timezone date/time
          this.formData.selectedStartTime = slotEl.dataset.start;
          this.formData.selectedEndTime = slotEl.dataset.end;
          this.formData.selectedDate = slotEl.dataset.date; // Practice timezone
          this.formData.selectedTime = slotEl.dataset.time; // Practice timezone
          this.formData.clinicianId = slotEl.dataset.clinician || null;
          this.formData.displayTime = slotEl.dataset.display;
          
          console.log('HelpPet Widget: Selected slot', {
            startTime: this.formData.selectedStartTime,
            endTime: this.formData.selectedEndTime,
            date: this.formData.selectedDate,
            time: this.formData.selectedTime,
            clinician: this.formData.clinicianId,
            display: this.formData.displayTime
          });
          
          continueBtn.disabled = false;
        });
      });

      backBtn.addEventListener('click', () => this.handleBack());
      continueBtn.addEventListener('click', () => this.handleNext());
      
      // Date navigation buttons
      if (prevBtn && !prevBtn.disabled) {
        prevBtn.addEventListener('click', () => this.changeDate(-1));
      }
      if (nextDayBtn) {
        nextDayBtn.addEventListener('click', () => this.changeDate(1));
      }
    }
  }

  // Global API - supports multiple instances
  window.HelpPetWidget = {
    VERSION: '1.0.0',
    
    create: function(config) {
      // Validate required fields
      if (!config.practiceId || !config.token || !config.containerId) {
        console.error('HelpPet Widget: Missing required configuration (practiceId, token, containerId)');
        return null;
      }
      
      // Clean inputs
      config.token = config.token.trim();
      config.practiceId = config.practiceId.trim();
      
      // Check for placeholder/obfuscated tokens
      if (config.token.includes('•') || config.token === 'YOUR_API_TOKEN' || config.token.length < 10) {
        console.error('HelpPet Widget: Invalid token - appears to be placeholder or obfuscated');
        console.error('HelpPet Widget: Please copy the real token from your dashboard');
        return null;
      }
      
      console.log('HelpPet Widget v' + this.VERSION + ': Creating instance');
      
      return new HelpPetWidget(config);
    },
    
    // Backwards compatibility
    init: function(config) {
      return this.create(config);
    }
  };
})();


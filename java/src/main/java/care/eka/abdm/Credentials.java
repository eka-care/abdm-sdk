package care.eka.abdm;

import java.time.Duration;
import java.time.Instant;

/**
 * Authentication credentials for API access. Mirrors {@code auth.Credentials}
 * in Go and {@code Credentials} in Python's credentials module.
 */
public record Credentials(
        String accessToken,
        String refreshToken,
        Instant expiresAt,
        Instant refreshExpiresAt,
        String source) {

    /** Matches Go's Credentials.Expired 5-minute buffer. */
    private static final Duration EXPIRY_BUFFER = Duration.ofMinutes(5);

    /** Returns true if the access token is expired (with the 5-minute buffer). */
    public boolean expired() {
        return Instant.now().isAfter(expiresAt.minus(EXPIRY_BUFFER));
    }

    /** Returns true if the refresh token is present and not yet expired. */
    public boolean canRefresh() {
        return refreshToken != null && !refreshToken.isEmpty()
                && Instant.now().isBefore(refreshExpiresAt);
    }
}

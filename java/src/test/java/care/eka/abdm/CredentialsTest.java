package care.eka.abdm;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import org.junit.jupiter.api.Test;

class CredentialsTest {

    @Test
    void notExpiredWellBeforeExpiry() {
        Credentials creds = new Credentials("access", "refresh",
                Instant.now().plus(1, ChronoUnit.HOURS), Instant.now().plus(2, ChronoUnit.HOURS), "test");
        assertFalse(creds.expired());
    }

    @Test
    void expiredInsideTheFiveMinuteBuffer() {
        Credentials creds = new Credentials("access", "refresh",
                Instant.now().plus(2, ChronoUnit.MINUTES), Instant.now().plus(1, ChronoUnit.HOURS), "test");
        assertTrue(creds.expired(), "a token expiring in 2 minutes must be treated as expired (5-minute buffer)");
    }

    @Test
    void expiredAfterActualExpiry() {
        Credentials creds = new Credentials("access", "refresh",
                Instant.now().minus(1, ChronoUnit.MINUTES), Instant.now().plus(1, ChronoUnit.HOURS), "test");
        assertTrue(creds.expired());
    }

    @Test
    void canRefreshWithUsableRefreshToken() {
        Credentials creds = new Credentials("access", "refresh",
                Instant.now().minus(1, ChronoUnit.MINUTES), Instant.now().plus(1, ChronoUnit.HOURS), "test");
        assertTrue(creds.canRefresh());
    }

    @Test
    void cannotRefreshWithEmptyRefreshToken() {
        Credentials creds = new Credentials("access", "",
                Instant.now().minus(1, ChronoUnit.MINUTES), Instant.now().plus(1, ChronoUnit.HOURS), "test");
        assertFalse(creds.canRefresh());
    }

    @Test
    void cannotRefreshOnceRefreshTokenExpired() {
        Credentials creds = new Credentials("access", "refresh",
                Instant.now().minus(1, ChronoUnit.MINUTES), Instant.now().minus(1, ChronoUnit.SECONDS), "test");
        assertFalse(creds.canRefresh());
    }
}
